package volumes

import (
	"fmt"
	"strings"
)

const (
	SELECT_LANGUAGE  = "SELECT_LANGUAGE"
	LOGIN            = "LOGIN"
	PASSWORD         = "PASSWORD"
	END_CONVERSATION = "END_CONVERSATION"
	CHANGE_LANGUAGE  = "CHANGE_LANGUAGE"
	SUBMIT_NAME      = "SUBMIT_NAME"
	SUBMIT_PHONE     = "SUBMIT_PHONE"
)

type UserSession struct {
	State    string
	Language string
	Username string
	Name     string
	Phone    string
	Password string
}

const RemoteServerURL = "http://84.46.247.18/api/v1/internet-tariffs/public?offset=0&limit=100"

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
