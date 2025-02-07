package chat

import (
	"fmt"
	"strings"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var cachedCategories []volumes.CategoryDataType       // Assuming the type returned by server.GetCategories is server.Category
var cachedSubCategories []volumes.SubCategoryDataType // Assuming the type returned by server.GetCategories is server.Category

func handleCategorySelect(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID

	selectedCategoryName := update.Message.Text

	var selectedCategoryID int64
	found := false
	// Find the ID of the category based on its name
	for _, category := range cachedCategories {
		if strings.TrimSpace(category.Name) == strings.TrimSpace(selectedCategoryName) {
			selectedCategoryID = category.Id
			found = true
			break
		}
	}
	if !found {
		return
	}

	lang := "uz"
	token := ""

	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.SelectedCategoryId = selectedCategoryID
		lang = user.Language
		token = user.Token
	}
	var err error
	// If there's a valid token, fetch the user balance

	cachedSubCategories, err = server.GetSubCategories(lang, token, selectedCategoryID, -1)

	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Error fetching data from the server.")
		bot.Send(msg)
		return
	}

	// Create a new keyboard with category buttons
	var keyboard [][]tgbotapi.KeyboardButton

	// Map cachedSubCategories to keyboard buttons
	for _, subCategory := range cachedSubCategories {
		button := tgbotapi.NewKeyboardButton(subCategory.Question)
		keyboard = append(keyboard, []tgbotapi.KeyboardButton{button})
	}

	// Add the "main menu" button at the bottom
	mainMenuButton := tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu"))
	keyboard = append(keyboard, []tgbotapi.KeyboardButton{mainMenuButton})

	// Create the keyboard markup
	replyMarkup := tgbotapi.NewReplyKeyboard(keyboard...)
	replyMarkup.ResizeKeyboard = true

	// Send the message with the keyboard
	message := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "pleaseSelectFAQ"))
	message.ReplyMarkup = replyMarkup
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.State = volumes.SELECT_FAQ
	}
	bot.Send(message)
}

func handleSubCategorySelect(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID

	selectedFAQName := update.Message.Text

	var selectedSubCategoryID int64
	var selectedSubCategoryAnswer string

	// Find the ID of the category based on its name

	if strings.TrimSpace(selectedFAQName) == strings.TrimSpace(translations.GetTranslation(userSessions, chatID, "otherQuestion")) {
		//messageText := "Telegram: @turonsupport"
		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "operatorMessage"))
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return

	}

	for _, category := range cachedSubCategories {
		if strings.TrimSpace(category.Question) == strings.TrimSpace(selectedFAQName) {
			selectedSubCategoryID = category.Id
			selectedSubCategoryAnswer = category.Answer
			break
		}
	}

	lang := "uz"
	token := ""
	var selectedCategoryID int64

	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.SelectedSubCategoryId = selectedSubCategoryID
		lang = user.Language
		selectedCategoryID = user.SelectedCategoryId
		token = user.Token
	}
	var err error
	// If there's a valid token, fetch the user balance

	cachedSubCategories, err = server.GetSubCategories(lang, token, selectedCategoryID, selectedSubCategoryID)

	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Error fetching data from the server.")
		bot.Send(msg)
		return
	}

	// Create a new keyboard with category buttons
	var keyboard [][]tgbotapi.KeyboardButton

	// Map cachedSubCategories to keyboard buttons
	if len(cachedSubCategories) > 0 {
		for _, category := range cachedSubCategories {
			button := tgbotapi.NewKeyboardButton(category.Question)
			keyboard = append(keyboard, []tgbotapi.KeyboardButton{button})
		}
	}

	// Add the "main menu" button at the bottom
	mainMenuButton := tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu"))
	connectToButton := tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "otherQuestion"))
	keyboard = append(keyboard, []tgbotapi.KeyboardButton{connectToButton})
	keyboard = append(keyboard, []tgbotapi.KeyboardButton{mainMenuButton})
	// Create the keyboard markup
	replyMarkup := tgbotapi.NewReplyKeyboard(keyboard...)
	replyMarkup.ResizeKeyboard = true

	var message tgbotapi.MessageConfig

	if selectedSubCategoryAnswer != "" {
		message = tgbotapi.NewMessage(chatID, selectedSubCategoryAnswer)

	} else {
		message = tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "pleaseSelectFAQ"))
	}

	// Send the message with the keyboard
	message.ReplyMarkup = replyMarkup
	bot.Send(message)
}

func ShowCategories(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	// Check if the user session exists
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)

		// If there's no token, change the user state to LOGIN
		if user.Phone == "" {
			user.State = volumes.SUBMIT_PHONE
			contactButton := tgbotapi.NewKeyboardButton(fmt.Sprintf("📱 %s", translations.GetTranslation(userSessions, chatID, "sharePhoneNumber")))
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

			msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "pleaseShareYourPhoneNumber"))
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			return
		}

		// If there's no token, change the user state to LOGIN
		if user.Token == "" {
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
		}
		var err error
		cachedCategories, err = server.GetCategories(user.Language)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Error fetching data from the server."))
			return
		}

		// Create a new keyboard with category buttons
		var keyboard [][]tgbotapi.KeyboardButton

		// Map cachedCategories to keyboard buttons
		for _, category := range cachedCategories {
			button := tgbotapi.NewKeyboardButton(category.Name)
			keyboard = append(keyboard, []tgbotapi.KeyboardButton{button})

		}
		// Add the "main menu" button at the bottom
		mainMenuButton := tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu"))
		keyboard = append(keyboard, []tgbotapi.KeyboardButton{mainMenuButton})

		// Create the keyboard markup
		replyMarkup := tgbotapi.NewReplyKeyboard(keyboard...)
		replyMarkup.ResizeKeyboard = true

		// Send the message with the keyboard
		message := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "pleaseSelectCategory"))
		message.ReplyMarkup = replyMarkup
		if session, ok := userSessions.Load(chatID); ok {
			user := session.(*volumes.UserSession)
			user.State = volumes.SELECT_CATEGORY
		}
		bot.Send(message)
	}
}

func HandleChatConversation(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map, user *volumes.UserSession) {

	switch user.State {
	case volumes.SELECT_CATEGORY:
		handleCategorySelect(bot, update, userSessions)
		user.State = volumes.SELECT_FAQ

	case volumes.SELECT_FAQ:
		handleSubCategorySelect(bot, update, userSessions)
		user.State = volumes.SELECT_FAQ
	}
}
