package conversations

import (
	"fmt"
	"log"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/handlers/events"
	"github.com/OzodbekX/TuronMiniApp/helpers"
	"github.com/OzodbekX/TuronMiniApp/server"

	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

func handleLogOut(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	switch update.Message.Text {
	case translations.GetTranslation(userSessions, chatID, "yes"):
		if session, ok := userSessions.Load(chatID); ok {
			user := session.(*volumes.UserSession)
			errorResponse := server.TerminateOwnSession(user)
			fmt.Println(errorResponse)
		}
		helpers.StartEvent(bot, chatID, userSessions)
	case translations.GetTranslation(userSessions, chatID, "no"):
		events.ShowMainMenu(bot, chatID, userSessions)
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
	loginRespose, err := server.LoginToBackend(user.Phone, user.Username, password, userID)
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
	user.Token = loginRespose.AccessToken
	user.RefreshToken = loginRespose.RefreshToken
	user.TuronId = loginRespose.TuronId
	// Assuming `balanceData` is fetched and has the required fields
	balanceData, err := server.GetUserData(user)
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
	msg.ParseMode = "HTML"
	deleteUserMessage(bot, chatID, update.Message.MessageID)

	bot.Send(msg)
	events.ShowMainMenu(bot, chatID, userSessions)
	user.State = volumes.END_CONVERSATION

}
