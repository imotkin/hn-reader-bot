package main

import (
	"log"
	"os"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	bot := NewBot(token, true)

	if err := bot.Init(); err != nil {
		log.Fatalf("initialize bot error: %v", err)
	}

	bot.Start()
}
