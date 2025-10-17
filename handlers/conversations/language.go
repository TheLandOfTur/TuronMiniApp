package conversations

import (
	"fmt"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/handlers/events"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleLanguage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	lang := "uz"

	switch update.Message.Text {
	case "\U0001F1F7\U0001F1FA Русский":
		lang = "ru"
	case "\U0001F1FA\U0001F1FF O'zbekcha":
		lang = "uz"
	}
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.Language = lang
		user.State = volumes.SUBMIT_PHONE
	}
	contactButton := tgbotapi.NewKeyboardButton(fmt.Sprintf("📱 %s", translations.GetTranslation(userSessions, chatID, "sharePhoneNumber")))
	contactButton.RequestContact = true // Enable the contact request

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			contactButton,
		),
		// tgbotapi.NewKeyboardButtonRow(
		// 	tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu")),
		// ),
	)
	keyboard.OneTimeKeyboard = true // Show keyboard only once
	keyboard.ResizeKeyboard = true  // Adjust keyboard size to fit the screen

	msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "pleaseShareYourPhoneNumber"))
	msg.ReplyMarkup = keyboard

	// langKeyboard := tgbotapi.NewReplyKeyboard(
	// 	tgbotapi.NewKeyboardButtonRow(
	// 		tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "cancel")),
	// 		tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu")),
	// 	),
	// )
	// msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "login"))
	// msg.ReplyMarkup = langKeyboard
	bot.Send(msg)
}
func onchangeLanguage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	lang := "uz"
	switch update.Message.Text {
	case "\U0001F1F7\U0001F1FA Русский":
		lang = "ru"
	case "\U0001F1FA\U0001F1FF O'zbekcha":
		lang = "uz"
	}
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.Language = lang
		user.State = volumes.END_CONVERSATION
	}

	events.ShowMainMenu(bot, chatID, userSessions)
}
