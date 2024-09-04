package main

import (
	"cmp"
	"flag"
	"log"
	"os"
)

func main() {
	tokenFlag := flag.String("token", "", "The token to connect Telegram API")
	flag.Parse()

	token := cmp.Or(*tokenFlag, os.Getenv("BOT_TOKEN"))
	bot := NewBot(token, true)

	if err := bot.Init(); err != nil {
		log.Fatalf("initialize bot error: %v", err)
	}

	bot.Start()
}
