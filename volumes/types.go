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
	RefreshToken          string
	TuronId               int64
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

// Create a struct for the request payload
type UserSessionRequest struct {
	OSName string `json:"OSName"`
	//OSVersion   string `json:"OSVersion"`
	//DeviceModel string `json:"deviceModel"`
	//DeviceName  string `json:"deviceName"`
	DeviceType string `json:"deviceType"`
}
type PromoCodeResponse struct {
	Status  string      `json:"status"`  // "OK", "ALREADY_EXISTS", "PERMISSION_DENIED"
	Data    interface{} `json:"data"`    // Flexible to handle string or nested data
	Success bool        `json:"success"` // true or false
}

// SpeedByTime represents the time-based speed limits
type SpeedByTime struct {
	FromTime int `json:"fromTime"`
	Speed    int `json:"speed"`
	ToTime   int `json:"toTime"`
}

// BillingTariffPlan represents the main plan structure
type BillingTariffPlan struct {
	ApartmentTypes   *[]string      `json:"apartmentTypes"`
	BillingTariffID  int            `json:"billingTariffId"`
	BillingType      string         `json:"billingType"`
	CostPerByte      int            `json:"costPerByte"`
	CostPerMb        int            `json:"costPerMb"`
	CostPeriod       int            `json:"costPeriod"`
	CreatedAt        string         `json:"createdAt"` // Consider time.Time if dealing with timestamps
	Description      *string        `json:"description"`
	DevicesMaxAmount int            `json:"devicesMaxAmount"`
	DevicesMinAmount int            `json:"devicesMinAmount"`
	ID               string         `json:"id"`
	Image            string         `json:"image"`
	ImageMobile      string         `json:"imageMobile"`
	IsActive         bool           `json:"isActive"`
	Name             string         `json:"name"`
	Position         int            `json:"position"`
	PrepaidTraffic   int            `json:"prepaidTraffic"`
	Price            int            `json:"price"`
	SpeedByTime      *[]SpeedByTime `json:"speedByTime"` // Struct for handling time-based speed
	TrafficInet      int            `json:"trafficInet"`
	Type             string         `json:"type"`
	UpdatedAt        string         `json:"updatedAt"` // Consider time.Time if handling timestamps
}
type BalanceData struct {
	DateStart                string            `json:"dateStart"`
	EndDate                  string            `json:"endDate"`
	Address                  string            `json:"address"`
	Apartment                string            `json:"apartment"`
	Identify                 string            `json:"identify"`
	Login                    string            `json:"login"`
	Phone                    string            `json:"phone"`
	Balance                  int               `json:"balance"`
	DiscountLoyality         float64           `json:"discountLoyality"`
	AbonentId                int               `json:"abonentId"`
	AdditionalTraffic        int               `json:"additionalTraffic"`
	UnreadNotificationsCount int               `json:"unreadNotificationsCount"`
	Tariff                   BillingTariffPlan `json:"tariff"`
	//NextTariff               *BillingTariffPlan `json:"nextTariff"`
	//RecommendedTariff        *BillingTariffPlan `json:"recommendedTariff"`
}
type SubCategoryResponse struct {
	Success bool                  `json:"success"`
	Data    []SubCategoryDataType `json:"data"`
}

type TariffSpeedObject struct {
	FromTime int `json:"fromTime"`
	Speed    int `json:"speed"`
	ToTime   int `json:"toTime"`
}
type TariffObject struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Price       int                 `json:"price"`
	SpeedByTime []TariffSpeedObject `json:"speedByTime"`
}
type MetaObject struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
type ResponseTpe struct {
	Data    []TariffObject `json:"data"`
	Meta    MetaObject     `json:"meta"`
	Status  string         `json:"status"`
	Success bool           `json:"success"`
}

type CategoryResponse struct {
	Success bool               `json:"success"`
	Data    []CategoryDataType `json:"data"`
}

// Create a struct for the request payload
type LoginRequest struct {
	Login          string             `json:"login"`
	Password       string             `json:"password"`
	PhoneNumber    string             `json:"phoneNumber"`
	TelegramUserID string             `json:"telegramUserId"`
	UserSession    UserSessionRequest `json:"userSession"`
}

// Create a struct for the response
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TuronId      int64  `json:"turonId"`
}

type LoginResponse struct {
	Status  string        `json:"status"`
	Success bool          `json:"success"`
	Data    TokenResponse `json:"data"`
}

type SubscriptionResponse struct {
	Status  string      `json:"status"`
	Success bool        `json:"success"`
	Data    BalanceData `json:"data"`
}
