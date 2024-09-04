package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	messageTemplate      = "<b><a href=\"https://news.ycombinator.com/item?id=%d\">%s</a></b>\n<a href=\"https://news.ycombinator.com/user?id=%s\">%s</a> / %s / Comments: %d\n\n"
	messageNotFound      = "The command isn't found!"
	messageInvalidFormat = "Your message doesn't contain any command. Please try again!"
	messageError         = "Server error... oops"

	loggingTemplate = "Time: %s / User: %s / Text: %s"

	storiesLimit = 10

	timeFormatFull  = "02-01-2006 15:04:05"
	timeFormatShort = "02-01-2006 15:04"
)

var (
	sb strings.Builder

	ErrEmptyToken = errors.New("empty token")
	ErrConnectAPI = errors.New("telegram API connection")

	commands = []tgbotapi.BotCommand{
		{Command: "/start", Description: "Start the bot"},
		{Command: "/best", Description: "Show best stories"},
		{Command: "/new", Description: "Show new stories"},
		{Command: "/top", Description: "Show top stories"},
		{Command: "/job", Description: "Show job stories"},
		{Command: "/ask", Description: "Show ask stories"},
		{Command: "/poll", Description: "Show poll stories"},
	}
)

type Bot struct {
	api     *tgbotapi.BotAPI
	client  *ClientAPI
	token   string
	logging bool
}

func NewBot(token string, logging bool) *Bot {
	return &Bot{
		token:   token,
		logging: logging,
		client:  NewClientAPI(),
	}
}

func (b *Bot) Init() error {
	if b.token == "" {
		return ErrEmptyToken
	}

	api, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return ErrConnectAPI
	}

	b.api = api

	cmds := tgbotapi.NewSetMyCommands(commands...)

	if _, err = b.api.Request(cmds); err != nil {
		return fmt.Errorf("set commands: %w", err)
	}

	return nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if b.logging {
			b.LogUpdate(update)
		}

		msg := b.HandleMessage(update.Message)

		_, err := b.api.Send(msg)
		if err != nil {
			log.Printf("send message: %v", err)
		}
	}
}

func (b *Bot) HandleMessage(msg *tgbotapi.Message) (response tgbotapi.MessageConfig) {
	if msg == nil {
		return
	}

	response = tgbotapi.NewMessage(msg.Chat.ID, "")

	if !msg.IsCommand() {
		response.Text = messageInvalidFormat
		return
	}

	command := StoriesType(msg.Command())

	switch command {
	case New, Best, Top, Ask, Job, Poll:
		stories, err := b.client.Stories(command)
		if err != nil {
			response.Text = messageError
			return
		}

		response.Text = b.PrintStories(stories)
		response.ParseMode = tgbotapi.ModeHTML
	default:
		response.Text = messageNotFound
	}

	return
}

func (b *Bot) PrintStories(stories []int) string {
	sb.Reset()

	selected := stories[:storiesLimit]

	for _, id := range selected {
		story, _ := b.client.GetStory(id)
		sb.WriteString(fmt.Sprintf(
			messageTemplate,
			story.ID, story.Title,
			story.By, story.By,
			time.Unix(story.Time, 0).Format(timeFormatShort),
			story.Descendants,
		))
	}

	return sb.String()
}

func (b *Bot) LogUpdate(u tgbotapi.Update) {
	log.Printf(
		loggingTemplate,
		u.Message.Time().Format(timeFormatFull),
		u.SentFrom().UserName,
		u.Message.Text,
	)
}
