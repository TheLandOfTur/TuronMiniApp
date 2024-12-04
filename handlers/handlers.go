package handlers

import (
	"fmt"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/handlers/conversations"
	"github.com/OzodbekX/TuronMiniApp/handlers/events"
	"github.com/OzodbekX/TuronMiniApp/translations"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserSession struct {
	State    string
	Language string
	Username string
	Password string
}

var userSessions = sync.Map{}

func HandleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := update.Message
	chatID := msg.Chat.ID

	if msg.Text == "/start" {
		conversations.StartEvent(bot, chatID, &userSessions)
		return
	}
	if msg.Text == fmt.Sprintf("🚪 %s", translations.GetTranslation(&userSessions, chatID, "Exit")) {
		conversations.StartEvent(bot, chatID, &userSessions)
		return
	}
	if msg.Text == translations.GetTranslation(&userSessions, chatID, "cancel") {
		conversations.StartEvent(bot, chatID, &userSessions)
		return
	}
	if msg.Text == translations.GetTranslation(&userSessions, chatID, "mainMenu") {
		events.ShowMainMenu(bot, chatID, &userSessions)
		return
	}
	if msg.Text == fmt.Sprintf("🌐 %s", translations.GetTranslation(&userSessions, chatID, "Language")) {
		events.ShowLanguages(bot, chatID, &userSessions)
		return
	}
	if msg.Text == fmt.Sprintf("📝 %s", translations.GetTranslation(&userSessions, chatID, "Application")) {
		events.SendRequestToBackend(bot, chatID, &userSessions)
		return
	}

	conversations.HandleUpdateConversation(bot, update, &userSessions)

}
