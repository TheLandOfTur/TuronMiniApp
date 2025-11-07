package helpers

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/joho/godotenv"
)

func cutFirst16Chars(dateStr string) string {
	parsedTime, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		return ""
	}

	// Format the date in "DD.MM.YYYY 00:00"
	return parsedTime.Format("02.01.2006 00:00")
}

func AddSpacesEveryThreeDigits(number int) string {
	numStr := fmt.Sprintf("%d", number) // Convert the number to a string
	var result strings.Builder
	if -999 < number && number < 999 {
		return numStr
	}

	// Iterate over the string in reverse
	if number < 0 {
		numStr = strings.ReplaceAll(numStr, "-", "")
	}
	length := len(numStr)
	for i, ch := range numStr {
		if (length-i)%3 == 0 && i != 0 { // Add a space every 3 digits, but not at the start
			result.WriteRune(' ')
		}
		result.WriteRune(ch)
	}
	if number < 0 {
		return fmt.Sprintf("- %s", result.String())
	}

	return result.String()
}
func ConvertDateFormat(input string) string {
	// Define the input and output formats explicitly
	inputFormat := "2006-01-02"
	outputFormat := "02.01.2006"

	// Parse the input date string into a time.Time object
	parsedDate, err := time.Parse(inputFormat, input)
	if err != nil {
		return ""
	}

	// Format the parsed date into the desired format
	formattedDate := parsedDate.Format(outputFormat)
	return formattedDate
}

// Get formatted subscription message
func GetSubscriptionMessage(balanceData volumes.BalanceData, chatID int64, userSessions *sync.Map) (string, error) {
	// Get translations based on the user language
	translate := func(key string) string {
		return translations.GetTranslation(userSessions, chatID, key)
	}

	translateDate := func() string {
		lang := "uz"
		if session, ok := userSessions.Load(chatID); ok {
			user := session.(*volumes.UserSession)
			lang = user.Language
		}
		if lang == "uz" {
			return ConvertDateFormat(balanceData.DateStart) + " " + translate("from") + " " + ConvertDateFormat(balanceData.EndDate) + " " + translate("to")

		} else {
			return translate("from") + " " + ConvertDateFormat(balanceData.DateStart) + " " + translate("to") + " " + ConvertDateFormat(balanceData.EndDate)
		}
	}
	subscriptionStatus := "active" // default to active
	subscriptionIcon := "\U0001F7E2"

	if balanceData.Balance < 0 {
		subscriptionStatus = "inactive"
		subscriptionIcon = "\U0001F534"
	}

	// Create the message with translated fields
	formattedMessage := fmt.Sprintf(
		"<b>%s</b>: %s %s\n"+
			"<b>%s</b>: %s\n"+
			"<b>%s</b>: %s %s\n"+
			"<b>%s</b>: %s\n"+
			"<b>%s</b>: %s\n"+
			"<b>%s</b>: %s%s",
		translate("yourBalance"), // Translated "Your current balance"
		AddSpacesEveryThreeDigits(balanceData.Balance),
		translate("uzs"), // uzs

		translate("tariffName"), // Translated "Tariff Name"
		balanceData.Tariff.Name,

		translate("subscriptionPrice"), // Translated "Subscription Price"
		AddSpacesEveryThreeDigits(int(balanceData.Tariff.Price)),
		translate("uzs"),

		translate("nextSubscriptionDate"), // Translated "Next Subscription Date"
		cutFirst16Chars(balanceData.EndDate),

		translate("subscriptionPeriod"), // Translated "Subscription Period"
		translateDate(),
		translate("subscriptionActive"),                 // Translated "Subscription Active"
		translate(subscriptionStatus), subscriptionIcon, // Translated "Active"/"Inactive"
	)

	return formattedMessage, nil
}

func MustToken() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	TOKEN := os.Getenv("TOKEN")
	if TOKEN == "" {
		log.Fatalf("TOKEN is not set in .env file")
	}
	return TOKEN
}

func GetLastMessageID(bot *tgbotapi.BotAPI, chatID int64) (int, error) {
	// Fetch updates from Telegram
	updates, err := bot.GetUpdates(tgbotapi.UpdateConfig{
		Offset:  0,   // Start from the latest unprocessed update
		Limit:   100, // Max number of updates to fetch
		Timeout: 0,   // No long polling
	})
	if err != nil {
		return 0, fmt.Errorf("error fetching updates: %v", err)
	}

	// Iterate through updates and find the last message for the given chatID
	var lastMessageID int
	for _, update := range updates {
		if update.Message != nil && update.Message.Chat.ID == chatID {
			lastMessageID = update.Message.MessageID
		}
	}

	if lastMessageID == 0 {
		return 0, fmt.Errorf("no messages found for chatID: %d", chatID)
	}

	return lastMessageID, nil
}

// GetFormattedPromoCodeMessage generates a user-friendly message based on the promo code response
func GetFormattedPromoCodeMessage(promoResponse volumes.PromoCodeResponse, chatID int64, userSessions *sync.Map) (string, error) {
	// Get translations based on the user language
	translate := func(key string) string {
		return translations.GetTranslation(userSessions, chatID, key)
	}
	// Default icon and status based on the response status
	statusIcon := "\U0001F534"                      // Default to red (failure)
	statusMessage := translate("promoCodeNotFound") // Default to "Promo code inactive"

	if promoResponse.Status == "OK" && promoResponse.Success {
		statusIcon = "\U0001F7E2"                    // Green circle for success
		statusMessage = translate("promoCodeActive") // "Promo code active"
	} else if promoResponse.Status == "ALREADY_EXISTS" {
		statusIcon = "\U0001F7E1"                              // Yellow circle for warning
		statusMessage = translate("promoCodeAlreadyActivated") // "Promo code already activated"
	} else if promoResponse.Status == "PERMISSION_DENIED" {
		statusIcon = "\U0001F6AB"                              // Prohibited symbol for access denied
		statusMessage = translate("promoCodePermissionDenied") // "Permission denied when entering promo code"
	}

	// Generate the formatted message
	formattedMessage := fmt.Sprintf(
		"<b>%s</b>: %s %s\n",
		translate("status"), // Translated "Promo Code Status"
		statusMessage,
		statusIcon,
	)

	return formattedMessage, nil
}

// getUserAgent generates a User-Agent string based on the platform.
func GetUserSession() volumes.UserSessionRequest {
	// Get the operating system and architecture
	os := runtime.GOOS
	//arch := runtime.GOARCH

	// Define a base User-Agent string for your bot
	baseUserAgent := "BOT"

	// Append platform-specific information
	switch os {
	case "windows":
		return volumes.UserSessionRequest{
			OSName:     "Windows",
			DeviceType: baseUserAgent,
		}
	case "darwin":
		return volumes.UserSessionRequest{
			OSName:     "macOS",
			DeviceType: baseUserAgent,
		}
	case "linux":
		return volumes.UserSessionRequest{
			OSName:     "Linux",
			DeviceType: baseUserAgent,
		}
	default:
		return volumes.UserSessionRequest{
			OSName:     os,
			DeviceType: baseUserAgent,
		}
	}

}

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

func GetBaseFAQUrl(apiPath string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	BASE_FAQ_URL := os.Getenv("BASE_FAQ_URL")
	if BASE_FAQ_URL == "" {
		log.Fatalf("BASE_FAQ_URL is not set in .env file")
	}
	url := fmt.Sprintf("%s%s", BASE_FAQ_URL, apiPath)
	return url
}

func DeleteTemporaryMessages(bot *tgbotapi.BotAPI, chatID int64, user *volumes.UserSession) {
	if len(user.TemporaryMessages) == 0 {
		return // nothing to delete
	}

	for _, msgID := range user.TemporaryMessages {
		if chatID == 0 || msgID <= 0 {
			continue // skip invalid IDs silently
		}

		deleteMsg := tgbotapi.NewDeleteMessage(chatID, int(msgID))
		if _, err := bot.Request(deleteMsg); err != nil {
			// Log the error, but continue with the rest
			log.Printf("[WARN] Could not delete message %d: %v", msgID, err)
			continue
		}
	}
	// Clear the slice after deletion
	user.TemporaryMessages = []int{}
}
