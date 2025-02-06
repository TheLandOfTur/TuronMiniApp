package volumes

const (
	SELECT_LANGUAGE    = "SELECT_LANGUAGE"
	LOGIN              = "LOGIN"
	PASSWORD           = "PASSWORD"
	END_CONVERSATION   = "END_CONVERSATION"
	CHANGE_LANGUAGE    = "CHANGE_LANGUAGE"
	SUBMIT_NAME        = "SUBMIT_NAME"
	SUBMIT_PHONE       = "SUBMIT_PHONE"
	SELECT_CATEGORY    = "SELECT_CATEGORY"
	SELECT_FAQ         = "SELECT_FAQ"
	LOG_OUT            = "LOG_OUT"
	ACTIVATE_PROMOCODE = "ACTIVATE_PROMOCODE"
)

type UserSession struct {
	State                 string
	Language              string
	Username              string
	Name                  string
	Phone                 string
	Password              string
	Token                 string
	SelectedCategoryId    int64
	SelectedSubCategoryId int64
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
type CategoryDataType struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SubCategoryDataType struct {
	Id       int64  `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer,omitempty"`
}
