package conversations

import (
	"regexp"
	"strings"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func isValidPhoneNumber(phoneNumber string) bool {
	regex := regexp.MustCompile(`^(?:\+998|\b998)?\d{9}$`)
	return regex.MatchString(phoneNumber)
}

func handlePhoneNumber(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID

	// Extract phone number from Contact or Text
	var phoneNumber string
	if update.Message.Contact != nil {
		phoneNumber = update.Message.Contact.PhoneNumber // Shared via contact button
	} else {
		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "sharePhoneNumber"))
		bot.Send(msg)
		return
		// 		phoneNumber = update.Message.Text // User manually enters the phone number
	}

	// Validate phone number format
	if !isValidPhoneNumber(phoneNumber) {
		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "sharePhoneNumber"))
		bot.Send(msg)
		return
	}

	// Update the user's session if the number is valid
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		if !strings.HasPrefix(phoneNumber, "+") {
			user.Phone = "+" + phoneNumber
		} else {
			user.Phone = phoneNumber
		}
		user.State = volumes.CHOOSE_USER_TYPE
	}

	langKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "abonent")),
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "user")),
		),
	)
	msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "chooseRole"))
	msg.ReplyMarkup = langKeyboard
	bot.Send(msg)
}
