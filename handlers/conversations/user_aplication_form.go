package conversations

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/helpers"
	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleFullNameInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	msg := update.Message
	if msg == nil {
		log.Println("[WARN] HandleFullNameInput called with nil message")
		return
	}
	chatID := msg.Chat.ID
	text := update.Message.Text
	if text == translations.GetTranslation(userSessions, chatID, "abonent") || text == translations.GetTranslation(userSessions, chatID, "user") {
		identifyUserType(bot, update, userSessions)
		return
	}
	text = strings.TrimSpace(text)

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
		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		retryMsg.ReplyMarkup = removeKeyboard
		bot.Send(retryMsg)
		return
	}

	// ✅ Save full name to user session
	user.FullName = text
	user.State = volumes.CHOOSE_REGIONS
	fetchRegions(bot, chatID, userSessions)

}

func writeLocationItems(bot *tgbotapi.BotAPI, user *volumes.UserSession, chatID int64, regions []volumes.Region, title string, isRegion bool) {
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
	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Printf("[ERROR] Failed to send message: %v", err)
		return
	}

	// Assuming user is your *volumes.UserSession
	user.TemporaryMessages = append(user.TemporaryMessages, sentMsg.MessageID)
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
		writeLocationItems(bot, user, chatID, regions, translations.GetTranslation(userSessions, chatID, "pleaseSelectYurDistrict"), true)
	}

}

func handleRegionWrite(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	if text == translations.GetTranslation(userSessions, chatID, "abonent") || text == translations.GetTranslation(userSessions, chatID, "user") {
		identifyUserType(bot, update, userSessions)
		return
	}
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		if text == translations.GetTranslation(userSessions, chatID, "backOneStep") {
			msg := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "chooseRole"))
			roleKeyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "abonent")),
					tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "user")),
				),
			)
			msg.ReplyMarkup = roleKeyboard
			if _, err := bot.Send(msg); err != nil {
				log.Printf("[ERROR] Failed to send full name prompt: %v", err)
			}
			user.State = volumes.CHOOSE_USER_TYPE
			return
		}
		helpers.DeleteTemporaryMessages(bot, chatID, user)
		fetchRegions(bot, chatID, userSessions)
	}
}

func handeDistrictWrite(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	chatID := update.Message.Chat.ID
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		text := update.Message.Text
		if text == translations.GetTranslation(userSessions, chatID, "backOneStep") {
			helpers.DeleteTemporaryMessages(bot, chatID, user)
			fetchRegions(bot, chatID, userSessions)
			user.State = volumes.CHOOSE_REGIONS
			return
		}
		summary := createApplicationText(userSessions, user, chatID, 0)
		// 🧩 Send confirmation message
		msgToSend := tgbotapi.NewMessage(chatID, summary)
		backKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "backOneStep")),
			),
		)
		msgToSend.ReplyMarkup = backKeyboard
		sentMsg, err := bot.Send(msgToSend)
		if err != nil {
			log.Printf("[ERROR] Failed to send confirmation message: %v", err)
		}
		user.TemporaryMessages = append(user.TemporaryMessages, sentMsg.MessageID)

		districts, err := server.GetDistricts(user, int64(user.RegionId))
		if err != nil {
			log.Printf("⚠️ Error fetching districts")
			return
		}
		user.State = volumes.CHOOSE_DISTRICTS
		user.Districts = districts
		writeLocationItems(bot, user, chatID, districts, translations.GetTranslation(userSessions, chatID, "pleaseSelectYurRegion"), false)
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
		user.State = volumes.CHOOSE_DISTRICTS
		summary := createApplicationText(userSessions, user, chatID, 0)
		// 🧩 Send confirmation message
		msgToSend := tgbotapi.NewMessage(chatID, summary)
		backKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "backOneStep")),
			),
		)
		msgToSend.ReplyMarkup = backKeyboard
		sentMsg, err := bot.Send(msgToSend)
		if err != nil {
			log.Printf("[ERROR] Failed to send full name prompt: %v", err)
		}
		user.TemporaryMessages = append(user.TemporaryMessages, sentMsg.MessageID)

		districts, err := server.GetDistricts(user, int64(regionID))
		if err != nil {
			bot.Request(tgbotapi.NewCallback(callback.ID, "⚠️ Error fetching districts"))
			return
		}
		user.Districts = districts
		writeLocationItems(bot, user, chatID, districts, translations.GetTranslation(userSessions, chatID, "pleaseSelectYurRegion"), false)

	}
}

func HandleDistrictSelection(bot *tgbotapi.BotAPI, update *tgbotapi.Update, userSessions *sync.Map) {
	callback := update.CallbackQuery
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	data := callback.Data
	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		helpers.DeleteTemporaryMessages(bot, chatID, user)
	}

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
		user.State = volumes.USER_CABINET
		telegramUserID := update.CallbackQuery.From.ID // 👈 this is the TelegramUserID
		var username string
		if callback.From != nil {
			username = callback.From.UserName
		}
		result, err := server.SendApplicationToBackend(
			user.RegionId,
			getRegionName(user),
			user.DistrictId,
			getDistrictName(user),
			"user",
			user.Phone,
			user.Language,
			telegramUserID,
			username,
		)

		if err != nil {
			log.Println("❌ Error sending application:", err)
			return
		}
		var requestId int64 = 0
		if result.Data != nil {
			// Try to interpret Data as a map (most JSON API responses decode into map[string]interface{})
			if dataMap, ok := result.Data.(map[string]interface{}); ok {
				if idVal, ok := dataMap["id"]; ok {
					// Handle if backend sends numeric ID (float64 is default for JSON numbers)
					switch v := idVal.(type) {
					case float64:
						requestId = int64(v)
					case int:
						requestId = int64(v)
					case int64:
						requestId = v
					}
				}
			}
		}
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(translations.GetTranslation(userSessions, chatID, "connectWithOperator"), "https://t.me/TuronTelecomSales"),
			),
		)
		var applicationText = createApplicationText(userSessions, user, chatID, requestId)
		message := tgbotapi.NewMessage(chatID, applicationText)
		message.ReplyMarkup = inlineKeyboard
		sentMsg, err := bot.Send(message)
		if err != nil {
			log.Printf("[ERROR] Failed to send phone number prompt: %v", err)
		}
		if err != nil {
			log.Printf("[ERROR] Failed to send confirmation message: %v", err)
		}
		user.TemporaryMessages = append(user.TemporaryMessages, sentMsg.MessageID)
		sendSuccessApplicationMessage(bot, user, userSessions, chatID)

	} else {
		log.Printf("[WARN] No session found for chatID: %d", chatID)
	}
}

func createApplicationText(userSessions *sync.Map, user *volumes.UserSession, chatID int64, requestId int64) string {
	var result strings.Builder

	// 1️⃣ Title (Application Ready)
	if requestId > 0 {
		result.WriteString(translations.GetTranslation(userSessions, chatID, "applicationSentSuccessfully") + "\n")
	}

	// 🔹 Translations
	//titleFullName := translations.GetTranslation(userSessions, chatID, "fullName") // 👤 Full name
	titlePhone := translations.GetTranslation(userSessions, chatID, "phoneNumber") // 📞 Phone
	titleRegion := translations.GetTranslation(userSessions, chatID, "city")       // 🏙️ Region
	titleDistrict := translations.GetTranslation(userSessions, chatID, "district") // 📍 District

	// 2️⃣ Full Name
	//if strings.TrimSpace(user.FullName) != "" {
	//	result.WriteString(fmt.Sprintf("👤 %s: %s\n", titleFullName, user.FullName))
	//}

	// 3️⃣ Phone
	if strings.TrimSpace(user.Phone) != "" {
		result.WriteString(fmt.Sprintf("📞 %s: %s\n", titlePhone, user.Phone))
	}

	// 4️⃣ Region
	var regionName string
	for _, r := range user.Regions {
		if r.ID == user.RegionId {
			regionName = r.Name
			break
		}
	}
	if strings.TrimSpace(regionName) != "" {
		result.WriteString(fmt.Sprintf("🏙️ %s: %s\n", titleRegion, regionName))
	}

	// 5️⃣ District
	var districtName string
	for _, d := range user.Districts {
		if d.ID == user.DistrictId {
			districtName = d.Name
			break
		}
	}
	if strings.TrimSpace(districtName) != "" {
		result.WriteString(fmt.Sprintf("📍 %s: %s\n", titleDistrict, districtName))
	}

	// Remove any trailing newline
	return strings.TrimSpace(result.String())
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
			tgbotapi.NewKeyboardButton(translations.GetTranslation(userSessions, chatID, "exit1")),
		),
	)
	selectUserTypeMessage := tgbotapi.NewMessage(chatID, translations.GetTranslation(userSessions, chatID, "PleaseSelectOption"))
	selectUserTypeMessage.ReplyMarkup = langKeyboard
	bot.Send(selectUserTypeMessage)
	user.State = volumes.USER_CABINET
	return
}
