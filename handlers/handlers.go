package handlers

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/helpers"

	"github.com/OzodbekX/TuronMiniApp/handlers/chat"

	"github.com/OzodbekX/TuronMiniApp/handlers/conversations"
	"github.com/OzodbekX/TuronMiniApp/handlers/events"
	"github.com/OzodbekX/TuronMiniApp/translations"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userSessions = sync.Map{}

func HandleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := update.Message
	chatID := msg.Chat.ID

	if msg.Text == "/start" {
		helpers.StartEvent(bot, chatID, &userSessions)
		return
	}
	if msg.Text == fmt.Sprintf("🚪 %s", translations.GetTranslation(&userSessions, chatID, "Exit")) {
		//helpers.StartEvent(bot, chatID, &userSessions)
		events.QuestionaryLogOut(bot, chatID, &userSessions)
		return
	}
	if msg.Text == fmt.Sprintf("⬅️ %s", translations.GetTranslation(&userSessions, chatID, "GoBack")) {
		//helpers.StartEvent(bot, chatID, &userSessions)
		events.OnClickGoBack(bot, chatID, &userSessions)
		return
	}
	if msg.Text == translations.GetTranslation(&userSessions, chatID, "cancel") {
		helpers.StartEvent(bot, chatID, &userSessions)
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
	if msg.Text == fmt.Sprintf("💰 %s", translations.GetTranslation(&userSessions, chatID, "Balance")) {
		events.ShowUserBalance(bot, chatID, &userSessions)
		return
	}
	if msg.Text == fmt.Sprintf("🏷️ %s", translations.GetTranslation(&userSessions, chatID, "promoCode")) {
		events.RedirectToPromoCode(bot, chatID, &userSessions)
		return
	}
	if msg.Text == fmt.Sprintf("📊 %s", translations.GetTranslation(&userSessions, chatID, "Tariffs")) {
		events.ShowTariffList(bot, chatID, &userSessions)
		return
	}
	if msg.Text == fmt.Sprintf("❓ %s", translations.GetTranslation(&userSessions, chatID, "FAQ")) {
		//messageText := "Telegram: @turonsupport"
		//events.SendMessage(bot, chatID, messageText)
		chat.ShowCategories(bot, chatID, &userSessions)
		return
	}

	conversations.HandleUpdateConversation(bot, update, &userSessions)
}

func HandleInlineTaps(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] panic recovered in handleInlineTaps: %v\n%s", r, debug.Stack())
			if update != nil && update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
				chatID := update.CallbackQuery.Message.Chat.ID
				msg := tgbotapi.NewMessage(chatID, "⚠️ Something went wrong. Please try again.")
				_, _ = bot.Send(msg)
			}
		}
	}()

	// --- Validate the callback query ---
	if update.CallbackQuery == nil {
		log.Println("[WARN] handleInlineTaps called with nil CallbackQuery")
		return
	}

	callback := update.CallbackQuery
	data := callback.Data
	chatID := callback.Message.Chat.ID

	log.Printf("[DEBUG] Inline tap detected from chatID=%d, data=%s", chatID, data)

	// Always acknowledge callback (removes Telegram's "loading" spinner)
	_, _ = bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	// --- Route callback based on prefix ---
	switch {
	case strings.HasPrefix(data, "district_"):
		conversations.HandleDistrictSelection(bot, update, &userSessions)
	case strings.HasPrefix(data, "region_"):
		print("333333333333333333")
		conversations.HandleRegionSelection(bot, update, &userSessions)

	default:
		log.Printf("[WARN] Unknown callback data: %s", data)
		msg := tgbotapi.NewMessage(chatID, "❓ Unknown action.")
		_, _ = bot.Send(msg)
	}
}
