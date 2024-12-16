package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/OzodbekX/TuronMiniApp/volumes"
	"github.com/joho/godotenv"
)

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

	var objects ResponseTpe

	err = json.NewDecoder(resp.Body).Decode(&objects)

	if err != nil {
		return nil, err
	}

	return objects.Data, nil
}

type BalanceData struct {
	Balance              int    `json:"balance"`
	TariffName           string `json:"tariffName"`
	SubscriptionPrice    int    `json:"subscriptionPrice"`
	NextSubscriptionDate string `json:"nextSubscriptionDate"`
	StartPeriodDate      string `json:"startPeriodDate"`
	EndPeriodDate        string `json:"endPeriodDate"`
	SubscriptionStatus   bool   `json:"subscriptionStatus"`
	TuronID              int    `json:"turonId"`
}

type SubscriptionResponse struct {
	Status  string      `json:"status"`
	Success bool        `json:"success"`
	Data    BalanceData `json:"data"`
}

// GetUserData fetches user data from the server
func GetUserData(token string, language string) (BalanceData, error) {
	url := getBaseUrl("/api/v1/users/info")

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

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return BalanceData{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

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
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Perform the request
	resp, err := client.Do(req)
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

	var emptyArray = []volumes.CategoryDataType{}

	// Perform the request
	resp, err := client.Do(req)
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

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return emptyArray, fmt.Errorf("request failed: %w", err)
	}
	fmt.Println(resp)

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
