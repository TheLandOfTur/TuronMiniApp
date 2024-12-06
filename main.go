package main

import (
	"log"

	"github.com/OzodbekX/TuronMiniApp/handlers" // Import the handlers package
	"github.com/OzodbekX/TuronMiniApp/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := helpers.MustToken()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Main update loop
	for update := range updates {
		if update.Message != nil { // Check if update contains a message
			handlers.HandleMessage(bot, &update)
		}
	}
}
