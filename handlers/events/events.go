package events

import (
	"fmt"

	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func OnSubmitLogin(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user any) {
	chatID := msg.Chat.ID
	lang := user.(string)
	if len(msg.Text) > 0 && len(msg.Text) <= 50 && contains(msg.Text, ":") {
		// Simulate a login process
		credentials := msg.Text
		data, err := server.GetUserDataFromServer(credentials)
		if err != nil {
			reply := tgbotapi.NewMessage(chatID, "Вход не выполнен. / Kirish muvaffaqiyatsiz.")
			bot.Send(reply)
		} else {
			reply := tgbotapi.NewMessage(chatID, fmt.Sprintf("Привет %s, Email: %s / Salom %s, Email: %s", data.Name, data.Email, data.Name, data.Email))
			bot.Send(reply)
		}
		return
	} else {
		reply := tgbotapi.NewMessage(chatID, translations.GetTranslation(lang, "login"))
		bot.Send(reply)
	}
}
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
