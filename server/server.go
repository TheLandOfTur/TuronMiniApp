package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OzodbekX/TuronMiniApp/logger"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/OzodbekX/TuronMiniApp/volumes"
	"github.com/joho/godotenv"
)

var loggers = logger.GetLogger()

type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getBaseUrl(apiPath string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatalf("BASE_URL is not set in .env file")
	}
	url := fmt.Sprintf("%s%s", baseURL, apiPath)
	return url
}
func getBaseFAQUrl(apiPath string) string {
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

// Fetch objects from a server
func FetchTariffsFromServer() ([]TariffObject, error) {
	url := getBaseUrl("/api/v1/internet-tariffs/public?offset=0&limit=100")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	req.Header.Add("Language", `ru`)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	loggers.Info("response from tariffs", url, " ", resp)
	var objects ResponseTpe

	err = json.NewDecoder(resp.Body).Decode(&objects)

	if err != nil {
		return nil, err
	}

	return objects.Data, nil
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
	DiscountLoyality         int               `json:"discountLoyality"`
	AbonentId                int               `json:"abonentId"`
	AdditionalTraffic        int               `json:"additionalTraffic"`
	UnreadNotificationsCount int               `json:"unreadNotificationsCount"`
	Tariff                   BillingTariffPlan `json:"tariff"`
	//NextTariff               *BillingTariffPlan `json:"nextTariff"`
	//RecommendedTariff        *BillingTariffPlan `json:"recommendedTariff"`
}

type SubscriptionResponse struct {
	Status  string      `json:"status"`
	Success bool        `json:"success"`
	Data    BalanceData `json:"data"`
}
type PromoCodeResponse struct {
	Status  string      `json:"status"`  // "OK", "ALREADY_EXISTS", "PERMISSION_DENIED"
	Data    interface{} `json:"data"`    // Flexible to handle string or nested data
	Success bool        `json:"success"` // true or false
}

// GetUserData fetches user data from the server
func GetUserData(token string, language string) (BalanceData, error) {
	url := getBaseUrl("/api/v1/abonents/info")

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BalanceData{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Add("Language", language)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	loggers.Info("response from get user data", req)

	// Perform the request
	resp, err := client.Do(req)

	if err != nil {
		return BalanceData{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	loggers.Info("response from get user data", resp)

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read the body for additional error context
		return BalanceData{}, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode the response
	var subscriptionResponse SubscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&subscriptionResponse); err != nil {
		return BalanceData{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validate the response status
	if subscriptionResponse.Status != "OK" || !subscriptionResponse.Success {
		return BalanceData{}, fmt.Errorf("unsuccessful response: status = %s, success = %v", subscriptionResponse.Status, subscriptionResponse.Success)
	}

	// Return the data
	return subscriptionResponse.Data, nil
}

// Submit user token
func ActivateToken(token string, pinCode string) (PromoCodeResponse, error) {
	url := getBaseUrl("/api/v1/abonents/activate-promo-code")

	// Create HTTP client and request
	client := &http.Client{}
	type PinCode struct {
		PinCode string `json:"pinCode"`
	}
	// Build the request payload
	payload := PinCode{
		PinCode: pinCode,
	}
	jsonPayload, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return PromoCodeResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	loggers.Info("response from activate function", url, " ", req)

	// Perform the request
	resp, err := client.Do(req)
	var promoCodeResponse PromoCodeResponse
	loggers.Info("response from activate function", resp)

	if err != nil {
		promoCodeResponse.Status = "UNKNOWN"
		return promoCodeResponse, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read the body for additional error context
		promoCodeResponse.Status = "UNKNOWN"
		promoCodeResponse.Success = false
		if resp.StatusCode == http.StatusForbidden {
			promoCodeResponse.Status = "PERMISSION_DENIED"
		}
		if resp.StatusCode == http.StatusConflict {
			promoCodeResponse.Status = "ALREADY_EXISTS"
		}
		return promoCodeResponse, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode the response
	if err := json.NewDecoder(resp.Body).Decode(&promoCodeResponse); err != nil {
		promoCodeResponse.Status = "UNKNOWN"

		return PromoCodeResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Return the data
	return promoCodeResponse, nil
}

// LoginToBackend logs in to the backend with phoneNumber, login, and password
func LoginToBackend(phoneNumber, login, password string, telegramUserID int64) (string, error) {
	url := getBaseUrl("/api/v1/users/sign-in-outside")

	// Create a struct for the request payload
	type LoginRequest struct {
		Login          string `json:"login"`
		Password       string `json:"password"`
		PhoneNumber    string `json:"phoneNumber"`
		TelegramUserID string `json:"telegramUserID"`
	}

	// Create a struct for the response
	type TokenResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	type LoginResponse struct {
		Status  string        `json:"status"`
		Success bool          `json:"success"`
		Data    TokenResponse `json:"data"`
	}

	// Build the request payload
	payload := LoginRequest{
		Login:          login,
		Password:       password,
		PhoneNumber:    phoneNumber,
		TelegramUserID: string(telegramUserID),
	}

	// Encode payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create an HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	loggers.Info("response from login to backend", url, " ", req)

	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Perform the request
	resp, err := client.Do(req)
	loggers.Info("response from login to backend", err, " ", resp)

	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 response codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode the JSON response
	var loginResponse LoginResponse
	if err := json.Unmarshal(body, &loginResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for a successful response
	if !loginResponse.Success || loginResponse.Status != "OK" {
		return "", fmt.Errorf("login failed: server response was not successful")
	}
	// Return the access token
	return loginResponse.Data.AccessToken, nil
}

type CategoryResponse struct {
	Success bool                       `json:"success"`
	Data    []volumes.CategoryDataType `json:"data"`
}

func GetCategories(language string) ([]volumes.CategoryDataType, error) {
	url := getBaseFAQUrl("/api/faqCategory/v1")

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		req.Header.Add("Accept", "/")
		return []volumes.CategoryDataType{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Add("Language", language)

	//req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	var emptyArray = []volumes.CategoryDataType{}

	// Perform the request
	loggers.Info("response from get categories", url, " ", req)

	resp, err := client.Do(req)
	if err != nil {
		return emptyArray, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()
	loggers.Info("response from get categories", err, " ", resp)

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read the body for additional error context

		return emptyArray, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}
	// Decode the response
	var subscriptionResponse CategoryResponse
	body, _ := io.ReadAll(resp.Body) // Read the body for additional error context

	if err := json.Unmarshal(body, &subscriptionResponse); err != nil {
		return emptyArray, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validate the response status
	if subscriptionResponse.Success != true || !subscriptionResponse.Success {
		return emptyArray, fmt.Errorf("unsuccessful response: status = %s, success = %v", "ok", subscriptionResponse.Success)
	}

	// Return the data
	return subscriptionResponse.Data, nil
}

type SubCategoryResponse struct {
	Success bool                          `json:"success"`
	Data    []volumes.SubCategoryDataType `json:"data"`
}

func GetSubCategories(lang, token string, categoryId, subCategoryId int64) ([]volumes.SubCategoryDataType, error) {
	var apiPath string
	if subCategoryId == -1 {
		apiPath = fmt.Sprintf("/api/faq/v1/withAnswer?categoryId=%d", categoryId)
	} else {
		apiPath = fmt.Sprintf("/api/faq/v1/withAnswer?categoryId=%d&parentFaqId=%d", categoryId, subCategoryId)

	}

	url := getBaseFAQUrl(apiPath)
	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		req.Header.Add("Accept", "/")
		return []volumes.SubCategoryDataType{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Add("Language", lang)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	var emptyArray = []volumes.SubCategoryDataType{}
	loggers.Info("response from get subcategories", url, " ", req)

	// Perform the request
	resp, err := client.Do(req)
	loggers.Info("response from subcategories", err, " ", resp)

	if err != nil {
		return emptyArray, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read the body for additional error context
		return emptyArray, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}
	// Decode the response
	var subscriptionResponse SubCategoryResponse
	body, _ := io.ReadAll(resp.Body) // Read the body for additional error context

	if err := json.Unmarshal(body, &subscriptionResponse); err != nil {
		return emptyArray, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validate the response status
	if subscriptionResponse.Success != true || !subscriptionResponse.Success {
		return emptyArray, fmt.Errorf("unsuccessful response: status = %s, success = %v", "ok", subscriptionResponse.Success)
	}
	// Return the data
	return subscriptionResponse.Data, nil
}
