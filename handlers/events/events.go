package events

import (
	"fmt"
	"log"
	"strings"
	"sync"

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
		messageBuilder.WriteString(fmt.Sprintf("%s: %s%s\n", translations.GetTranslation(userSessions, chatID, "price"), volumes.AddSpacesEveryThreeDigits(obj.Price), translations.GetTranslation(userSessions, chatID, "uzs")))
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
