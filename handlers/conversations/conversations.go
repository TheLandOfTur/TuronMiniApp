package conversations

import (
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/handlers/chat"

	"github.com/OzodbekX/TuronMiniApp/handlers/events"
	"github.com/OzodbekX/TuronMiniApp/helpers"
	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var lastMessageIDs sync.Map // To track the last message sent by the bot

func StartEvent(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {

	// Clear the user session
	userSessions.Delete(chatID)
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
	contactButton := tgbotapi.NewKeyboardButton(fmt.Sprintf("📱 %s", translations.GetTranslation(userSessions, chatID, "sharePhonenumber")))
	contactButton.RequestContact = true // Enable the contact request

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			contactButton,
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu")),
		),
	)
	keyboard.OneTimeKeyboard = true // Show keyboard only once
	keyboard.ResizeKeyboard = true  // Adjust keyboard size to fit the screen

	msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "enterPhone"))
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

func isValidPhoneNumber(phoneNumber string) bool {
	// Regex: starts with +998 followed by 9 digits OR just 9 digits
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

		phoneNumber = update.Message.Text // User manually enters the phone number
	}

	// Validate phone number format
	if !isValidPhoneNumber(phoneNumber) {
		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "invalidPhoneNumber"))
		bot.Send(msg)
		return
	}

	// Update the user's session if the number is valid
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.Phone = phoneNumber
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
	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send bot message: %v", err)
	} else {
		// Store the bot's message ID for future deletion
		lastMessageIDs.Store(chatID, sentMsg.MessageID)
	}

}

// DeleteUserMessage deletes the message sent by the user
func deleteUserMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := bot.Send(deleteMsg); err != nil {
		log.Printf("Failed to delete user message %d in chatID %d: %v", messageID, chatID, err)
	}
	// Retrieve and delete the bot's last message
	if botMessageID, ok := lastMessageIDs.Load(chatID); ok {
		deleteBotMsg := tgbotapi.NewDeleteMessage(chatID, botMessageID.(int))
		if _, err := bot.Send(deleteBotMsg); err != nil {
			log.Printf("Failed to delete bot message in chatID %d: %v", chatID, err)
		}
		// After deleting, remove it from the map
		lastMessageIDs.Delete(chatID)
	}
}

func handlePassword(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	password := update.Message.Text

	// Check if the user session exists
	session, ok := userSessions.Load(chatID)
	if !ok {
		deleteUserMessage(bot, chatID, update.Message.MessageID)
		msg := tgbotapi.NewMessage(chatID, "Session not found. Please start the login process again.")
		bot.Send(msg)
		sentMsg, err := bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send bot message: %v", err)
		} else {
			// Store the bot's message ID for future deletion
			lastMessageIDs.Store(chatID, sentMsg.MessageID)
		}
		return
	}

	user := session.(*volumes.UserSession)
	user.Password = password
	userID := update.Message.From.ID

	// Call backend login function
	token, err := server.LoginToBackend(user.Phone, user.Username, password, userID)
	if err != nil {
		// Login failed
		deleteUserMessage(bot, chatID, update.Message.MessageID)
		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "wrongParol"))
		// Reset to password state
		user.State = volumes.LOGIN
		sentMsg, err := bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send bot message: %v", err)
		} else {
			// Store the bot's message ID for future deletion
			lastMessageIDs.Store(chatID, sentMsg.MessageID)
		}
		return
	}

	// Save the token to the session if needed
	user.Token = token
	// Assuming `balanceData` is fetched and has the required fields
	balanceData, err := server.GetUserData(token, user.Language)
	if err != nil {
		deleteUserMessage(bot, chatID, update.Message.MessageID)
		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "wrongParol"))

		user.State = volumes.LOGIN

		sentMsg, err := bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send bot message: %v", err)
		} else {
			// Store the bot's message ID for future deletion
			lastMessageIDs.Store(chatID, sentMsg.MessageID)
		}
		return
	}
	// Get the formatted subscription message
	formattedMessage, err := helpers.GetSubscriptionMessage(balanceData, chatID, userSessions)
	if err != nil {
		deleteUserMessage(bot, chatID, update.Message.MessageID)

		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "wrongParol"))
		sentMsg, err := bot.Send(msg)
		user.State = volumes.LOGIN

		if err != nil {
			log.Printf("Failed to send bot message: %v", err)
		} else {
			// Store the bot's message ID for future deletion
			lastMessageIDs.Store(chatID, sentMsg.MessageID)
		}
		return
	}

	// Send the formatted message
	msg := tgbotapi.NewMessage(chatID, formattedMessage)
	deleteUserMessage(bot, chatID, update.Message.MessageID)

	bot.Send(msg)
	events.ShowMainMenu(bot, chatID, userSessions)
	user.State = volumes.END_CONVERSATION

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
	case volumes.SUBMIT_PHONE:
		handlePhoneNumber(bot, update, userSessions)
	case volumes.CHANGE_LANGUAGE:
		onchangeLanguage(bot, update, userSessions)
	case volumes.SELECT_CATEGORY, volumes.SELECT_FAQ:
		chat.HandleChatConversation(bot, update, userSessions, user)
	}
}
