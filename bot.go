package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	messageTemplate      = "<b><a href=\"https://news.ycombinator.com/item?id=%d\">%s</a></b>\n<a href=\"https://news.ycombinator.com/user?id=%s\">%s</a> / %s / Comments: %d\n\n"
	messageNotFound      = "The command isn't found!"
	messageInvalidFormat = "Your message doesn't contain any command. Please try again!"
	messageError         = "Server error... oops"

	storiesLimit = 10

	timeFormatFull  = "02-01-2006 15:04:05"
	timeFormatShort = "02-01-2006 15:04"
)

var text strings.Builder

type Bot struct {
	token string
}

func NewBot() *Bot {
	return &Bot{token: os.Getenv("BOT_TOKEN")}
}

func (b *Bot) Start() error {
	if b.token == "" {
		return errors.New("empty token")
	}

	bot, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return fmt.Errorf("telegram API connection: %w", err)
	}

	commandsList := tgbotapi.NewSetMyCommands([]tgbotapi.BotCommand{
		{Command: "/start", Description: "Start the bot"},
		{Command: "/best", Description: "Get best stories"},
		{Command: "/new", Description: "Get new stories"},
		{Command: "/top", Description: "Get top stories"}}...)

	if _, err = bot.Request(commandsList); err != nil {
		return fmt.Errorf("set commands: %w", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf(
			"Time: %s / User: %s / Text: %s",
			update.Message.Time().Format(timeFormatFull),
			update.SentFrom().UserName,
			update.Message.Text,
		)

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			if update.Message.IsCommand() {
				cmd := storiesType(update.Message.Command())

				switch cmd {
				case New, Best, Top:
					stories, err := Stories(cmd)
					if err != nil {
						msg.Text = messageError
						break
					}

					selected := stories[:storiesLimit]

					for _, id := range selected {
						story, _ := GetStory(id)
						text.WriteString(fmt.Sprintf(
							messageTemplate,
							story.ID, story.Title,
							story.By, story.By,
							time.Unix(story.Time, 0).Format(timeFormatShort),
							story.Descendants,
						))
					}

					msg.ParseMode = tgbotapi.ModeHTML
					msg.Text = text.String()
					text.Reset()
				default:
					msg.Text = messageNotFound
				}
			} else {
				msg.Text = messageInvalidFormat
			}

			_, err = bot.Send(msg)
			if err != nil {
				log.Printf("send message: %v", err)
			}
		}
	}

	return nil
}
