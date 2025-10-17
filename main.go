package main

import (
	"log"
	"net/http"

	"github.com/OzodbekX/TuronMiniApp/handlers" // Import the handlers package
	"github.com/OzodbekX/TuronMiniApp/helpers"
	"github.com/OzodbekX/TuronMiniApp/listeners"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := helpers.MustToken()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	// Register the requestHandler with the HTTP server
	http.HandleFunc("/api/send-messages", listeners.RequestHandler(bot))
	port := ":8080" // Change this to any available port
	log.Printf("Starting server on port %s", port)
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Main update loop
	for update := range updates {
		if update.CallbackQuery != nil {
			// User clicked a button
			handlers.HandleInlineTaps(bot, &update)
			continue
		}
		if update.Message != nil { // Check if update contains a message
			handlers.HandleMessage(bot, &update)
		}
	}
}
