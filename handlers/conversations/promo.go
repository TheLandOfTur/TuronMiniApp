package conversations

import (
	"fmt"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/handlers/events"
	"github.com/OzodbekX/TuronMiniApp/helpers"
	"github.com/OzodbekX/TuronMiniApp/server"

	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func checkActivePromoCode(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		promoResponse, err := server.ActivateToken(user, update.Message.Text)
		if err != nil {
			// Handle error fetching balance data
			msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "promoCodeNotFound"))
			bot.Send(msg)
			return
		}

		formattedMessage, err := helpers.GetFormattedPromoCodeMessage(promoResponse, chatID, userSessions)
		if err != nil {
			//Handle error formatting subscription data
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error formatting subscription data: %v", err))
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(chatID, formattedMessage)
		msg.ParseMode = "HTML"
		bot.Send(msg)
		// Change the user state to END_CONVERSATION after balance is shown
		user.State = volumes.END_CONVERSATION
		events.ShowMainMenu(bot, chatID, userSessions)

	} else {
		helpers.StartEvent(bot, chatID, userSessions)
	}
}
