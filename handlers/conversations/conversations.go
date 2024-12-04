package conversations

import (
	"sync"

	"github.com/OzodbekX/TuronMiniApp/handlers/events"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartEvent(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	// Language selection
	userSessions.Clear()
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
		user.State = volumes.LOGIN
	}
	bot.Send(reply)
}

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
		user.State = volumes.LOGIN
	}

	langKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "cancel")),
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu")),
		),
	)
	msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "login"))
	msg.ReplyMarkup = langKeyboard
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

func handleLogin(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	username := update.Message.Text

	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.Username = username
		user.State = volumes.PASSWORD
	}

	// Prompt for password
	msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "enterPassword"))
	bot.Send(msg)
}

func handlePassword(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	password := update.Message.Text

	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.Password = password
		user.State = volumes.END_CONVERSATION
	}

	// Validate credentials (placeholder logic)
	if password == "secret" { // Replace with actual validation
		msg := tgbotapi.NewMessage(chatID, "Login successful! Welcome to the main menu.")
		bot.Send(msg)

		// Show the main menu
		events.ShowMainMenu(bot, chatID, userSessions)
	} else {
		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "wrongParol"))
		bot.Send(msg)

		// Reset to password state
		if session, ok := userSessions.Load(chatID); ok {
			user := session.(*volumes.UserSession)
			user.State = volumes.LOGIN
		}
	}
}

func HandleUpdateConversation(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID

	session, _ := userSessions.LoadOrStore(chatID, &volumes.UserSession{State: volumes.SELECT_LANGUAGE})
	user := session.(*volumes.UserSession)
	switch user.State {
	case volumes.SELECT_LANGUAGE:
		handleLanguage(bot, update, userSessions)
	case volumes.LOGIN:
		handleLogin(bot, update, userSessions)
	case volumes.PASSWORD:
		handlePassword(bot, update, userSessions)
	case volumes.CHANGE_LANGUAGE:
		onchangeLanguage(bot, update, userSessions)
	case volumes.SUBMIT_NAME, volumes.SUBMIT_PHONE:
		HandleSubmissionConversation(bot, update, userSessions)
	}

}
