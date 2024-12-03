package handlers

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userSessions = sync.Map{}

func HandleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	if msg.Text == "/start" {
		// Language selection
		langKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Русский"),
				tgbotapi.NewKeyboardButton("O'zbekcha"),
			),
		)
		reply := tgbotapi.NewMessage(chatID, "Пожалуйста, выберите язык: / Iltimos, tilni tanlang:")
		reply.ReplyMarkup = langKeyboard
		bot.Send(reply)
		return
	}

	if msg.Text == "Русский" || msg.Text == "O'zbekcha" {
		lang := "ru"
		if msg.Text == "O'zbekcha" {
			lang = "uz"
		}
		userSessions.Store(chatID, lang)
		reply := tgbotapi.NewMessage(chatID, translations.getTranslation(lang, "login"))
		bot.Send(reply)
		return
	}

	// Check for login format (username:password)
	if user, ok := userSessions.Load(chatID); ok {
		lang := user.(string)
		if len(msg.Text) > 0 && len(msg.Text) <= 50 && contains(msg.Text, ":") {
			// Simulate a login process
			credentials := msg.Text
			data, err := server.getUserDataFromServer(credentials)
			if err != nil {
				reply := tgbotapi.NewMessage(chatID, "Вход не выполнен. / Kirish muvaffaqiyatsiz.")
				bot.Send(reply)
			} else {
				reply := tgbotapi.NewMessage(chatID, fmt.Sprintf("Привет %s, Email: %s / Salom %s, Email: %s", data.Name, data.Email, data.Name, data.Email))
				bot.Send(reply)
			}
			return
		} else {
			reply := tgbotapi.NewMessage(chatID, getTranslation(lang, "login"))
			bot.Send(reply)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
