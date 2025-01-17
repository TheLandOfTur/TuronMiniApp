package conversations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Struct for submission payload
type SubmissionPayload struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// Start the submission conversation
func StartSubmissionConversation(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	// Initialize session for submission
	session, _ := userSessions.LoadOrStore(chatID, &volumes.UserSession{State: volumes.SUBMIT_NAME})
	user := session.(*volumes.UserSession)
	user.State = volumes.SUBMIT_NAME

	// Prompt user for their name
	msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "enterName"))
	bot.Send(msg)
}

// Handle the submission process
func HandleSubmissionConversation(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	userInput := update.Message.Text

	session, _ := userSessions.LoadOrStore(chatID, &volumes.UserSession{})
	user := session.(*volumes.UserSession)

	switch user.State {
	case volumes.SUBMIT_NAME:
		// Save the name and ask for the phone number
		user.Name = userInput
		user.State = volumes.SUBMIT_PHONE
		contactButton := tgbotapi.NewKeyboardButton("📱 Share your phone number")
		contactButton.RequestContact = true // Enable the contact request
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				contactButton,
				tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu")),
			),
		)

		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "pleaseShareYourPhoneNumber"))
		msg.ReplyMarkup = keyboard
		bot.Send(msg)

	case volumes.SUBMIT_PHONE:
		// Save the phone number
		user.Phone = userInput
		user.State = volumes.END_CONVERSATION

		// Submit data to the server
		err := submitToServer(user.Name, user.Phone)
		if err != nil {
			log.Printf("Error submitting data: %v", err)
			msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "submissionFailed"))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "submissionSuccess"))
			bot.Send(msg)
		}

		// Reset the user's session state
		user.State = volumes.SELECT_LANGUAGE
	}
}

// Submit data to the server
func submitToServer(name, phone string) error {
	// Replace this with your actual server URL
	serverURL := "https://example.com/submit"

	// Create the payload
	payload := SubmissionPayload{Name: name, Phone: phone}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Make the POST request
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}
	return nil
}
