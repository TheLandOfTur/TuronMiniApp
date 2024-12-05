package volumes

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
	Token    string
}

const RemoteServerURL = "http://84.46.247.18/api/v1/internet-tariffs/public?offset=0&limit=100"
