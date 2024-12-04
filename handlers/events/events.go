package events

import (
	"fmt"
	"log"
	"sync"

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
			tgbotapi.NewKeyboardButton(fmt.Sprintf("📝 %s", translations.GetTranslation(userSessions, chatID, "Application"))),
			tgbotapi.NewKeyboardButton(fmt.Sprintf("🌐 %s", translations.GetTranslation(userSessions, chatID, "Language"))),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(fmt.Sprintf("💰 %s", translations.GetTranslation(userSessions, chatID, "Balance"))),
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
