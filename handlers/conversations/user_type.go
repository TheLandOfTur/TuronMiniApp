package conversations

import (
	"sync"

	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func identifyUserType(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// Handle user button choice first
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		switch text {
		case translations.GetTranslation(userSessions, chatID, "abonent"):
			// Abonent selected → proceed to LOGIN
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

		case translations.GetTranslation(userSessions, chatID, "user"):
			// Regular user selected → move to region selection
			user.State = volumes.CHOOSE_LOCATIONS
			// Call function to fetch regions (will implement later)
			fetchRegions(bot, chatID, userSessions)
			return
		}
	}

	// If not yet chosen abonent/user → ask user type
	roleKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "abonent")),
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "user")),
		),
	)
	msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "chooseRole"))
	msg.ReplyMarkup = roleKeyboard
	bot.Send(msg)
}
