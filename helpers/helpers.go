package helpers

import (
	"fmt"
	"strings"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/server"
	"github.com/OzodbekX/TuronMiniApp/translations"
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
	subscriptionStatus := "inactive" // default to inactive
	if balanceData.SubscriptionStatus {
		subscriptionStatus = "active"
	}

	// Create the message with translated fields
	formattedMessage := fmt.Sprintf(
		"%s: %s%s\n"+
			"%s: %s\n"+
			"%s: %d\n"+
			"%s: %s\n"+
			"%s: %s\n"+
			"%s: %s",
		translate("yourBalance"), // Translated "Your current balance"
		AddSpacesEveryThreeDigits(balanceData.Balance),
		translate("uzs"),        // uzs
		translate("tariffName"), // Translated "Tariff Name"
		balanceData.TariffName,
		translate("subscriptionPrice"), // Translated "Subscription Price"
		balanceData.SubscriptionPrice,
		translate("nextSubscriptionDate"), // Translated "Next Subscription Date"
		cutFirst16Chars(balanceData.NextSubscriptionDate),
		translate("subscriptionPeriod"), // Translated "Subscription Period"
		translate("from")+" "+balanceData.StartPeriodDate+" "+translate("to")+" "+balanceData.EndPeriodDate,
		translate("subscriptionActive"), // Translated "Subscription Active"
		translate(subscriptionStatus),   // Translated "Active"/"Inactive"

	)

	return formattedMessage, nil
}
