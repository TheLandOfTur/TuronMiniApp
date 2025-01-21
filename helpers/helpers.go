package helpers

import (
	"fmt"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/joho/godotenv"
)

func cutFirst16Chars(dateStr string) string {
	// Split the string by ":"
	parts := strings.Split(dateStr, ":")
	date := strings.Split(parts[0], " ")[0]
	hour := strings.Split(parts[0], " ")[1]

	// If there are at least two parts, join the first two
	if len(parts) >= 2 {
		return ConvertDateFormat(date) + " " + hour + ":" + parts[1]
	}

	// If there are less than two parts, return the string as is
	return dateStr
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
func GetSubscriptionMessage(balanceData server.BalanceData, chatID int64, userSessions *sync.Map) (string, error) {
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
			return ConvertDateFormat(balanceData.StartPeriodDate) + " " + translate("from") + " " + ConvertDateFormat(balanceData.EndPeriodDate) + " " + translate("to")

		} else {
			return translate("from") + " " + ConvertDateFormat(balanceData.StartPeriodDate) + " " + translate("to") + " " + ConvertDateFormat(balanceData.EndPeriodDate)
		}
	}
	subscriptionStatus := "inactive" // default to inactive
	subscriptionIcon := "\U0001F534"
	if balanceData.SubscriptionStatus {
		subscriptionStatus = "active"
		subscriptionIcon = "\U0001F7E2"
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
		translate("uzs"),        // uzs
		translate("tariffName"), // Translated "Tariff Name"
		balanceData.TariffName,
		translate("subscriptionPrice"), // Translated "Subscription Price"
		AddSpacesEveryThreeDigits(int(balanceData.SubscriptionPrice)),
		translate("uzs"),
		translate("nextSubscriptionDate"), // Translated "Next Subscription Date"
		cutFirst16Chars(balanceData.NextSubscriptionDate),
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
