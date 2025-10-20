package conversations

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func writeLocationItems(bot *tgbotapi.BotAPI, chatID int64, regions []volumes.Region, title string, isRegion bool) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var tempRow []tgbotapi.InlineKeyboardButton
	for i, region := range regions {
		var keyword = fmt.Sprintf("district_%d", region.ID)
		if isRegion {
			keyword = fmt.Sprintf("region_%d", region.ID)
		}

		tempRow = append(tempRow, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s", region.Name), // visible button text
			keyword,                        // callback data (contains ID)
		))
		// 2 buttons per line
		if (i+1)%2 == 0 {
			rows = append(rows, tempRow)
			tempRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	// Add the last row if uneven number
	if len(tempRow) > 0 {
		rows = append(rows, tempRow)
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, title)
	msg.ReplyMarkup = inlineKeyboard
	bot.Send(msg)
}

func fetchRegions(bot *tgbotapi.BotAPI, chatID int64, userSessions *sync.Map) {
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		regions, err := server.GetRegions(user)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Error fetching regions. Please try again later.")
			bot.Send(msg)
			return
		}
		user.Regions = regions
		writeLocationItems(bot, chatID, regions, translations.GetTranslation(userSessions, chatID, "pleaseSelectYurDistrict"), true)
	}

}

func HandleRegionSelection(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	callback := update.CallbackQuery
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	data := callback.Data
	_, _ = bot.Request(tgbotapi.NewCallback(callback.ID, ""))
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := bot.Request(deleteMsg); err != nil {
		log.Printf("[WARN] Failed to delete message %d: %v", messageID, err)
	}

	if !strings.HasPrefix(data, "region_") {
		return
	}

	// Extract the numeric ID
	idStr := strings.TrimPrefix(data, "region_")
	regionID, err := strconv.Atoi(idStr)
	if err != nil {
		bot.Request(tgbotapi.NewCallback(callback.ID, "⚠️ Invalid district"))
		return
	}

	// Save to user session
	if sessionData, ok := userSessions.Load(chatID); ok {
		user := sessionData.(*volumes.UserSession)
		user.RegionId = int64(regionID)
		user.State = volumes.CHOOSE_LOCATIONS
		districts, err := server.GetDistricts(user, int64(regionID))
		if err != nil {
			bot.Request(tgbotapi.NewCallback(callback.ID, "⚠️ Error fetching districts"))
			return
		}
		user.Districts = districts
		writeLocationItems(bot, chatID, districts, translations.GetTranslation(userSessions, chatID, "pleaseSelectYurRegion"), false)

	}
}

func HandleDistrictSelection(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	callback := update.CallbackQuery
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	data := callback.Data
	_, _ = bot.Request(tgbotapi.NewCallback(callback.ID, ""))
	if !strings.HasPrefix(data, "district_") {
		return
	}
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := bot.Request(deleteMsg); err != nil {
		log.Printf("[WARN] Failed to delete message %d: %v", messageID, err)
	}
	// Extract the numeric ID
	idStr := strings.TrimPrefix(data, "district_")
	districtID, err := strconv.Atoi(idStr)
	if err != nil {
		bot.Request(tgbotapi.NewCallback(callback.ID, "⚠️ Invalid district"))
		return
	}
	if sessionData, ok := userSessions.Load(chatID); ok {
		user := sessionData.(*volumes.UserSession)
		user.DistrictId = int64(districtID)
		user.State = volumes.ENTER_FULL_NAME
		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "enterFullName"))
		msg.ReplyMarkup = removeKeyboard
		if _, err := bot.Send(msg); err != nil {
			log.Printf("[ERROR] Failed to send full name prompt: %v", err)
		}
	} else {
		log.Printf("[WARN] No session found for chatID: %d", chatID)
	}
}

func createApplicationText(user *volumes.UserSession, additionalPhone ...string) string {
	// 🧩 Find selected Region and District names
	var regionName, districtName string

	for _, r := range user.Regions {
		if r.ID == user.RegionId {
			regionName = r.Name
			break
		}
	}

	for _, d := range user.Districts {
		if d.ID == user.DistrictId {
			districtName = d.Name
			break
		}
	}

	if regionName == "" {
		regionName = "—"
	}
	if districtName == "" {
		districtName = "—"
	}

	// 🧩 Build summary text
	fullNameText := fmt.Sprintf("👤 %s", user.FullName)
	regionText := fmt.Sprintf("🏙️ %s", regionName)
	districtText := fmt.Sprintf("📍 %s", districtName)

	phoneText := fmt.Sprintf("📞 %s", user.Phone)

	// 🧩 Handle optional additional phone
	additionalPhoneText := " "
	if len(additionalPhone) > 0 {
		additionalPhoneText = fmt.Sprintf("\n📞➕ %s", additionalPhone[0])
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s%s",
		fullNameText,
		regionText,
		districtText,
		phoneText,
		additionalPhoneText,
	)
}

func handleFullNameInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	msg := update.Message
	if msg == nil {
		log.Println("[WARN] HandleFullNameInput called with nil message")
		return
	}
	chatID := msg.Chat.ID
	text := strings.TrimSpace(msg.Text)

	sessionData, ok := userSessions.Load(chatID)
	if !ok {
		log.Printf("[WARN] No session found for chatID: %d", chatID)
		return
	}
	user := sessionData.(*volumes.UserSession)

	// 🧩 Validation: must be more than 3 letters
	if len([]rune(text)) < 3 {
		retryMsg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "fullNameTooShort"))
		if retryMsg.Text == "" {
			retryMsg.Text = translations.GetTranslation(userSessions, chatID, "placeEnterFullName")
		}
		bot.Send(retryMsg)
		return
	}

	// ✅ Save full name to user session
	user.FullName = text
	user.State = volumes.ENTER_ADDITIONAL_PHONE

	var applicationText = createApplicationText(user)
	message := tgbotapi.NewMessage(chatID, applicationText)
	if _, err := bot.Send(message); err != nil {
		log.Printf("[ERROR] Failed to send phone number prompt: %v", err)
	}
	// 🧩 Ask for additional phone number
	phonePrompt := translations.GetTranslation(userSessions, chatID, "enterAdditionalPhone")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "sendApplication")),
		),
	)
	keyboard.OneTimeKeyboard = true // Show keyboard only once
	keyboard.ResizeKeyboard = true  // Adjust keyboard size to fit the screen
	msgToSend := tgbotapi.NewMessage(chatID, phonePrompt)
	msgToSend.ReplyMarkup = keyboard

	if _, err := bot.Send(msgToSend); err != nil {
		log.Printf("[ERROR] Failed to send phone number prompt: %v", err)
	}
}

// getRegionName returns the name of the region selected by the user.
func getRegionName(user *volumes.UserSession) string {
	if user == nil {
		return ""
	}

	for _, region := range user.Regions {
		if region.ID == user.RegionId {
			return region.Name
		}
	}

	return ""
}

// getDistrictName returns the name of the district selected by the user.
func getDistrictName(user *volumes.UserSession) string {
	if user == nil {
		return ""
	}

	for _, district := range user.Districts {
		if district.ID == user.DistrictId {
			return district.Name
		}
	}

	return ""
}

func sendSuccessApplicationMessage(bot *tgbotapi.BotAPI, user *volumes.UserSession, userSessions *sync.Map, chatID int64) {
	// ✅ Success message
	langKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "myApplications")),
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "exit")),
		),
	)
	selectUserTypeMessage := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "applicationSentSuccessfully"))
	selectUserTypeMessage.ReplyMarkup = langKeyboard
	bot.Send(selectUserTypeMessage)
	user.State = volumes.SUCCESSFUL_STATE_USER
	return
}

func handleAdditionalPhoneInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	msg := update.Message
	if msg == nil {
		log.Println("[WARN] HandleAdditionalPhoneInput called with nil message")
		return
	}
	chatID := msg.Chat.ID
	text := strings.TrimSpace(msg.Text)
	telegramUserID := msg.From.ID // 👈 this is the TelegramUserID

	sessionData, ok := userSessions.Load(chatID)
	if !ok {
		log.Printf("[WARN] No session found for chatID: %d", chatID)
		return
	}
	user := sessionData.(*volumes.UserSession)

	sendAppText := translations.GetTranslation(userSessions, chatID, "sendApplication")
	// ✅ If user clicked “Send Application”
	if text == sendAppText {
		err := server.SendApplicationToBackend(
			user.RegionId,
			getRegionName(user),
			user.DistrictId,
			getDistrictName(user),
			user.FullName,
			user.Phone,
			user.Language,
			telegramUserID,
			nil,
		)

		if err != nil {
			log.Printf("[ERROR] Failed to send application: %v", err)
			bot.Send(tgbotapi.NewMessage(chatID, "❌ "+translations.GetTranslation(userSessions, chatID, "failedToSendApp")))
			return
		}

		sendSuccessApplicationMessage(bot, user, userSessions, chatID)
		return
	}
	// 🧩 Save phone number (optional)
	if text != "" {
		if !regexp.MustCompile(`^\+998\d{9}$`).MatchString(text) {
			errMsg := translations.GetTranslation(userSessions, chatID, "invalidPhoneNumber")
			if errMsg == "" {
				errMsg = translations.GetTranslation(userSessions, chatID, "wrongPhoneNumber")
			}
			bot.Send(tgbotapi.NewMessage(chatID, errMsg))
			return
		}
		user.AdditionalPhone = text
	}
	// 🧩 Move to next state
	user.State = volumes.CONFIRM_APPLICATION
	// 🧩 Add keyboard button to send application
	sendButton := tgbotapi.NewKeyboardButton(
		translations.GetTranslation(userSessions, chatID, "sendApplication"),
	)
	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(sendButton),
	)
	replyKeyboard.OneTimeKeyboard = true
	replyKeyboard.ResizeKeyboard = true
	summary := createApplicationText(user, text)
	// 🧩 Send confirmation message
	msgToSend := tgbotapi.NewMessage(chatID, summary)
	msgToSend.ReplyMarkup = replyKeyboard
	if _, err := bot.Send(msgToSend); err != nil {
		log.Printf("[ERROR] Failed to send confirmation message: %v", err)
	}
}

func handleFinalSubmission(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	msg := update.Message
	if msg == nil {
		log.Println("[WARN] HandleAdditionalPhoneInput called with nil message")
		return
	}
	chatID := msg.Chat.ID
	text := strings.TrimSpace(msg.Text)
	telegramUserID := msg.From.ID // 👈 this is the TelegramUserID

	sessionData, ok := userSessions.Load(chatID)
	if !ok {
		log.Printf("[WARN] No session found for chatID: %d", chatID)
		return
	}
	user := sessionData.(*volumes.UserSession)

	sendAppText := translations.GetTranslation(userSessions, chatID, "sendApplication")
	// ✅ If user clicked “Send Application”
	if text == sendAppText {
		err := server.SendApplicationToBackend(
			user.RegionId,
			getRegionName(user),
			user.DistrictId,
			getDistrictName(user),
			user.FullName,
			user.Phone,
			user.Language,
			telegramUserID,
			&user.AdditionalPhone,
		)
		if err != nil {
			log.Printf("[ERROR] Failed to send application: %v", err)
			bot.Send(tgbotapi.NewMessage(chatID, "❌ "+translations.GetTranslation(userSessions, chatID, "failedToSendApp")))
			return
		}

		sendSuccessApplicationMessage(bot, user, userSessions, chatID)
		return
	}

	user.State = volumes.CONFIRM_APPLICATION
	// 🧩 Add keyboard button to send application
	sendButton := tgbotapi.NewKeyboardButton(
		translations.GetTranslation(userSessions, chatID, "sendApplication"),
	)
	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(sendButton),
	)
	replyKeyboard.OneTimeKeyboard = true
	replyKeyboard.ResizeKeyboard = true
	summary := createApplicationText(user, user.AdditionalPhone)
	// 🧩 Send confirmation message
	msgToSend := tgbotapi.NewMessage(chatID, summary)
	msgToSend.ReplyMarkup = replyKeyboard
	if _, err := bot.Send(msgToSend); err != nil {
		log.Printf("[ERROR] Failed to send confirmation message: %v", err)
	}
}
