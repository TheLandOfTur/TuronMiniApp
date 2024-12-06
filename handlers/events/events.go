package events

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/helpers"
	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ShowMainMenu(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	// Create the keyboard for the main menu
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.State = volumes.END_CONVERSATION
	}
	mainMenuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(fmt.Sprintf("📊 %s", translations.GetTranslation(userSessions, chatID, "Tariffs"))),
			tgbotapi.NewKeyboardButton(fmt.Sprintf("❓ %s", translations.GetTranslation(userSessions, chatID, "FAQ"))),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(fmt.Sprintf("💰 %s", translations.GetTranslation(userSessions, chatID, "Balance"))),
			// tgbotapi.NewKeyboardButton(fmt.Sprintf("📝 %s", translations.GetTranslation(userSessions, chatID, "Application"))),
			tgbotapi.NewKeyboardButton(fmt.Sprintf("🌐 %s", translations.GetTranslation(userSessions, chatID, "Language"))),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(fmt.Sprintf("🚪 %s", translations.GetTranslation(userSessions, chatID, "Exit"))),
		),
	)
	// Create and send the message with the menu
	reply := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "PleaseSelectOption"))
	reply.ReplyMarkup = mainMenuKeyboard
	_, err := bot.Send(reply)
	if err != nil {
		// Handle error
		log.Printf("Error sending main menu: %v", err)
	}
}

func ShowUserBalance(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	// Check if the user session exists
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)

		// If there's no token, change the user state to LOGIN
		if user.Phone == "" {
			user.State = volumes.SUBMIT_PHONE
			contactButton := tgbotapi.NewKeyboardButton(fmt.Sprintf("📱 %s", translations.GetTranslation(userSessions, chatID, "sharePhonenumber")))
			contactButton.RequestContact = true // Enable the contact request

			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					contactButton,
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu")),
				),
			)
			keyboard.OneTimeKeyboard = true // Show keyboard only once
			keyboard.ResizeKeyboard = true  // Adjust keyboard size to fit the screen

			msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "enterPhone"))
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			return
		}

		// If there's no token, change the user state to LOGIN
		if user.Token == "" {
			user.State = volumes.LOGIN
			langKeyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "cancel")),
					tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu")),
				),
			)
			msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "login"))
			msg.ReplyMarkup = langKeyboard
			bot.Send(msg)
			return
		}

		// If there's a valid token, fetch the user balance
		balanceData, err := server.GetUserData(user.Token, user.Language)
		if err != nil {
			// Handle error fetching balance data
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch balance data: %v", err))
			bot.Send(msg)
			return
		}

		// Get the formatted subscription message
		formattedMessage, err := helpers.GetSubscriptionMessage(balanceData, chatID, userSessions)
		if err != nil {
			// Handle error formatting subscription data
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error formatting subscription data: %v", err))
			bot.Send(msg)
			return
		}

		// Send the formatted message
		msg := tgbotapi.NewMessage(chatID, formattedMessage)
		bot.Send(msg)

		// Change the user state to END_CONVERSATION after balance is shown
		user.State = volumes.END_CONVERSATION
	}
}

func ShowLanguages(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	// Language selection
	langKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("\U0001F1F7\U0001F1FA Русский"),
			tgbotapi.NewKeyboardButton("\U0001F1FA\U0001F1FF O'zbekcha"),
		),
	)
	reply := tgbotapi.NewMessage(chatID, "Пожалуйста, выберите язык: / Iltimos, tilni tanlang:")
	reply.ReplyMarkup = langKeyboard
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.State = volumes.CHANGE_LANGUAGE
	}
	bot.Send(reply)
}

func SendRequestToBackend(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	// Language selection
	langKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu")),
		),
	)
	reply := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "pleaseEnterYourName"))
	reply.ReplyMarkup = langKeyboard
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.State = volumes.SUBMIT_NAME
	}
	bot.Send(reply)
}

func ShowTariffList(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	// Replace with your server's endpoint

	// Fetch objects
	objects, err := server.FetchTariffsFromServer()
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Error fetching data from the server.")
		bot.Send(msg)
		return
	}

	// Helper function to convert seconds to hh:mm format
	convertSecondsToHHMM := func(seconds int) string {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		return fmt.Sprintf("%02d:%02d", hours, minutes) // Zero-padded format
	}

	// Prepare the message with all tariff data
	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("%s:\n\n", translations.GetTranslation(userSessions, chatID, "listOfTariffs")))

	for _, obj := range objects {
		messageBuilder.WriteString(fmt.Sprintf("<b>%s</b>\n", obj.Name))
		messageBuilder.WriteString(
			fmt.Sprintf("%s: %s%s\n", translations.GetTranslation(userSessions, chatID, "price"),
				helpers.AddSpacesEveryThreeDigits(obj.Price), translations.GetTranslation(userSessions, chatID, "uzs")),
		)
		messageBuilder.WriteString(fmt.Sprintf("%s:\n", translations.GetTranslation(userSessions, chatID, "speedByTime")))
		for _, speed := range obj.SpeedByTime {
			messageBuilder.WriteString(fmt.Sprintf("     %s - %s : %s %s \n",
				convertSecondsToHHMM(speed.FromTime),
				convertSecondsToHHMM(speed.ToTime),
				fmt.Sprintf("%d", speed.Speed/1000),
				translations.GetTranslation(userSessions, chatID, "mbs"),
			))
		}
		messageBuilder.WriteString("\n") // Add spacing between tariffs
	}

	// Send the formatted message
	msg := tgbotapi.NewMessage(chatID, messageBuilder.String())
	msg.ParseMode = "HTML"
	bot.Send(msg)
}
