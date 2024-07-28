package main

import "log"

func main() {
	bot := NewBot()
	if err := bot.Start(); err != nil {
		log.Fatalf("start bot error: %v", err)
	}
}
