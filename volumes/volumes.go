package volumes

const (
	SELECT_LANGUAGE  = "SELECT_LANGUAGE"
	LOGIN            = "LOGIN"
	PASSWORD         = "PASSWORD"
	END_CONVERSATION = "END_CONVERSATION"
	CHANGE_LANGUAGE  = "CHANGE_LANGUAGE"
	SUBMIT_NAME      = "SUBMIT_NAME"
	SUBMIT_PHONE     = "SUBMIT_PHONE"
	SELECT_CATEGORY  = "SELECT_CATEGORY"
	SELECT_SUBCAT    = "SELECT_SUBCAT"
	SELECT_FAQ       = "SELECT_FAQ"
)

type UserSession struct {
	State               string
	Language            string
	Username            string
	Name                string
	Phone               string
	Password            string
	Token               string
	SelectedCategory    int64
	SelectedSubCategory int64
}
type Message struct {
	Uz string `json:"uz,omitempty"` // Telegram user ID (optional)
	Ru string `json:"ru,omitempty"` // Message in Russian (optional)
	En string `json:"en,omitempty"` // Message in English (optional)
}

// RequestPayload represents the incoming HTTP request payload.
type AlertRequestPayload struct {
	Messages []struct {
		UserID  int64   `json:"userId"`  // Telegram user ID
		Message Message `json:"message"` // Message to send
	} `json:"messages"` // Array of user-message pairs
}
