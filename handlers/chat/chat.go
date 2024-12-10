package chat

import (
	"fmt"
	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"sync"
)

// UserSession holds the conversation state for a user
type UserSession struct {
	State            string
	SelectedCategory int
	SelectedSubCat   int
}

// Example data
var categories = []struct {
	ID   int
	Name string
}{
	{ID: 1, Name: "Internet"},
	{ID: 2, Name: "Mobile"},
}

var subCategories = []struct {
	ID       int
	Position int
	Question string
	Type     string
}{
	{ID: 1, Position: 1, Question: "Оплата и баланс", Type: "PASS_TO_DEFAULT"},
	{ID: 2, Position: 2, Question: "Как подключиться", Type: "FAQ"},
}

func handleCategorySelect(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID

	selectedCategory, err := strconv.ParseInt(update.Message.Text, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lang := "uz"

	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		user.SelectedCategory = selectedCategory
		user.State = volumes.SUBMIT_PHONE
		lang = user.Language

	}

	// If there's a valid token, fetch the user balance
	subCategories, err := server.GetSubCategories(lang, &selectedCategory, nil)
	fmt.Println("2222222222222222222")
	fmt.Println(subCategories)

	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Error fetching data from the server.")
		bot.Send(msg)
		return
	}
	fmt.Println(subCategories)

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

func sendSubCategories(bot *tgbotapi.BotAPI, chatID int64, categoryID int) {
	message := "Please select a subcategory:\n"
	for _, sub := range subCategories {
		message += sub.Question + "\n"
	}
	bot.Send(tgbotapi.NewMessage(chatID, message))
}

func sendFAQ(bot *tgbotapi.BotAPI, chatID int64, faq string) {
	bot.Send(tgbotapi.NewMessage(chatID, "FAQ: "+faq))
}

func getCategory(input string) *struct {
	ID   int
	Name string
} {
	for _, cat := range categories {
		if cat.Name == input {
			return &cat
		}
	}
	return nil
}

func getSubCategory(input string) *struct {
	ID       int
	Position int
	Question string
	Type     string
} {
	for _, sub := range subCategories {
		if sub.Question == input {
			return &sub
		}
	}
	return nil
}

func ShowCategories(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	// Check if the user session exists
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)

		// If there's no token, change the user state to LOGIN
		if user.Phone == "" {
			user.State = volumes.SUBMIT_PHONE
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

		// If there's a valid token, fetch the user balance
		categories, err := server.GetCategories(user.Language)

		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Error fetching data from the server.")
			bot.Send(msg)
			return
		}
		// Create a new keyboard with category buttons
		var keyboard [][]tgbotapi.KeyboardButton

		// Map categories to keyboard buttons
		fmt.Println(categories)
		var row []tgbotapi.KeyboardButton
		for _, category := range categories {
			button := tgbotapi.NewKeyboardButton(category.Name)
			row = append(row, button)
		}
		// Add the row of buttons to the keyboard
		keyboard = append(keyboard, row)

		// Add the "main menu" button at the bottom
		mainMenuButton := tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "mainMenu"))
		keyboard = append(keyboard, []tgbotapi.KeyboardButton{mainMenuButton})

		// Create the keyboard markup
		replyMarkup := tgbotapi.NewReplyKeyboard(keyboard...)

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
	chatID := update.Message.Chat.ID

	switch user.State {
	case volumes.SELECT_CATEGORY:
		handleCategorySelect(bot, update, userSessions)
		user.State = volumes.SELECT_SUBCAT

	case volumes.SELECT_SUBCAT:
		selectedCategory := getCategory(update.Message.Text)
		if selectedCategory != nil {
			//sendSubCategories(bot, chatID, user.SelectedCategory)
			user.State = volumes.SELECT_FAQ
		}

	case volumes.SELECT_FAQ:
		subCat := getSubCategory(update.Message.Text)
		if subCat != nil {
			if subCat.Type == "FAQ" {
				sendFAQ(bot, chatID, subCat.Question)
			} else {
				// Handle default action if not FAQ
				bot.Send(tgbotapi.NewMessage(chatID, "This subcategory doesn't have an FAQ. Moving to the default flow."))
			}
		}
	}
}
