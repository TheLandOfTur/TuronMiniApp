package volumes

const (
	SELECT_LANGUAGE  = "SELECT_LANGUAGE"
	LOGIN            = "LOGIN"
	PASSWORD         = "PASSWORD"
	END_CONVERSATION = "END_CONVERSATION"
)

type UserSession struct {
	State    string
	Language string
	Username string
	Password string
}
