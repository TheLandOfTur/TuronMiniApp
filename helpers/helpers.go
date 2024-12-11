package helpers

import (
	"fmt"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
	"github.com/joho/godotenv"
)

func cutFirst16Chars(dateStr string) string {
	// Split the string by ":"
	parts := strings.Split(dateStr, ":")

	// If there are at least two parts, join the first two
	if len(parts) >= 2 {
		return parts[0] + ":" + parts[1]
	}

	// If there are less than two parts, return the string as is
	return dateStr
}

func AddSpacesEveryThreeDigits(number int) string {
	numStr := fmt.Sprintf("%d", number) // Convert the number to a string
	var result strings.Builder
	if number < 999 {
		return numStr
	}

	// Iterate over the string in reverse
	length := len(numStr)
	for i, ch := range numStr {
		if (length-i)%3 == 0 && i != 0 { // Add a space every 3 digits, but not at the start
			result.WriteRune(' ')
		}
		result.WriteRune(ch)
	}

	return result.String()
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
			return balanceData.StartPeriodDate + " " + translate("from") + " " + balanceData.EndPeriodDate + " " + translate("to")

		} else {
			return translate("from") + " " + balanceData.StartPeriodDate + " " + translate("to") + " " + balanceData.EndPeriodDate
		}
	}
	subscriptionStatus := "inactive" // default to inactive
	if balanceData.SubscriptionStatus {
		subscriptionStatus = "active"
	}
	fmt.Printf("3333333333333333333333")
	fmt.Println(balanceData.SubscriptionPrice)

	// Create the message with translated fields
	formattedMessage := fmt.Sprintf(
		"%s: %s %s\n"+
			"%s: %s\n"+
			"%s: %s %s\n"+
			"%s: %s\n"+
			"%s: %s\n"+
			"%s: %s",
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
		translate("subscriptionActive"), // Translated "Subscription Active"
		translate(subscriptionStatus),   // Translated "Active"/"Inactive"
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
