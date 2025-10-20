package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/OzodbekX/TuronMiniApp/helpers"

	"github.com/OzodbekX/TuronMiniApp/logger"

	"github.com/OzodbekX/TuronMiniApp/volumes"
	"github.com/joho/godotenv"
)

var loggers = logger.GetLogger()

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

// RefreshToken refreshes the authentication token and executes the provided function on success
func RefreshToken[T any](onSuccess func() T, userData *volumes.UserSession) (T, error) {
	url := getBaseUrl("/api/v1/abonents/refresh-token")
	// Create an HTTP client and request
	client := &http.Client{}
	var zero T
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return zero, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userData.RefreshToken))
	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return zero, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 response codes
	if resp.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode the JSON response
	var refreshResponse = volumes.LoginResponse{}
	if err := json.Unmarshal(body, &refreshResponse); err != nil {
		return zero, fmt.Errorf("failed to decode response: %w", err)
	}

	if !refreshResponse.Success {
		return zero, errors.New("refresh token failed: server response was not successful")
	} else {
		userData.Token = refreshResponse.Data.AccessToken
		userData.RefreshToken = refreshResponse.Data.RefreshToken
		userData.TuronId = refreshResponse.Data.TuronId
	}

	return onSuccess(), nil
}

// Fetch objects from a server
func FetchTariffsFromServer() ([]volumes.TariffObject, error) {
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
	var objects volumes.ResponseTpe

	err = json.NewDecoder(resp.Body).Decode(&objects)

	if err != nil {
		return nil, err
	}

	return objects.Data, nil
}

// GetUserData fetches user data from the server
func GetUserData(userData *volumes.UserSession) (volumes.BalanceData, error) {
	url := getBaseUrl("/api/v1/abonents/info")
	var attemptRequest func() (volumes.BalanceData, error)
	attemptRequest = func() (volumes.BalanceData, error) {

		// Create HTTP client and request
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return volumes.BalanceData{}, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Add("Language", userData.Language)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userData.Token))
		loggers.Info("response from get user data", req)

		// Perform the request
		resp, err := client.Do(req)

		if err != nil {
			return volumes.BalanceData{}, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		loggers.Info("response from get user data", resp)

		// If server returns 401, refresh the token and retry
		if resp.StatusCode == http.StatusUnauthorized {
			loggers.Warn("received 401, attempting to refresh token")

			_, refreshErr := RefreshToken(func() volumes.TokenResponse {
				return volumes.TokenResponse{
					AccessToken:  userData.Token,
					TuronId:      userData.TuronId,
					RefreshToken: userData.RefreshToken,
				}
			}, userData)

			if refreshErr != nil {
				return volumes.BalanceData{}, fmt.Errorf("failed to refresh token: %w", refreshErr)
			}

			// Retry the request with the new token
			return attemptRequest()
		}

		// Check for non-200 status codes
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body) // Read the body for additional error context
			return volumes.BalanceData{}, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
		}

		// Decode the response
		var subscriptionResponse volumes.SubscriptionResponse
		if err := json.NewDecoder(resp.Body).Decode(&subscriptionResponse); err != nil {
			return volumes.BalanceData{}, fmt.Errorf("failed to decode response: %w", err)
		}

		// Validate the response status
		if subscriptionResponse.Status != "OK" || !subscriptionResponse.Success {
			return volumes.BalanceData{}, fmt.Errorf("unsuccessful response: status = %s, success = %v", subscriptionResponse.Status, subscriptionResponse.Success)
		}

		// Return the data
		return subscriptionResponse.Data, nil

	}
	return attemptRequest()

}

// Submit user token
func ActivateToken(userData *volumes.UserSession, pinCode string) (volumes.PromoCodeResponse, error) {
	url := getBaseUrl("/api/v1/abonents/activate-promo-code")
	var attemptRequest func() (volumes.PromoCodeResponse, error)
	attemptRequest = func() (volumes.PromoCodeResponse, error) {
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
			return volumes.PromoCodeResponse{}, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userData.Token))
		loggers.Info("response from activate function", url, " ", req)

		// Perform the request
		resp, err := client.Do(req)
		var promoCodeResponse volumes.PromoCodeResponse
		loggers.Info("response from activate function", resp)

		if err != nil {
			promoCodeResponse.Status = "UNKNOWN"
			return promoCodeResponse, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		loggers.Info("response from get user data", resp)

		// If server returns 401, refresh the token and retry
		if resp.StatusCode == http.StatusUnauthorized {
			loggers.Warn("received 401, attempting to refresh token")

			_, refreshErr := RefreshToken(func() volumes.TokenResponse {
				return volumes.TokenResponse{
					AccessToken:  userData.Token,
					TuronId:      userData.TuronId,
					RefreshToken: userData.RefreshToken,
				}
			}, userData)

			if refreshErr != nil {
				return volumes.PromoCodeResponse{}, fmt.Errorf("failed to refresh token: %w", refreshErr)
			}

			// Retry the request with the new token
			return attemptRequest()
		}

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

			return volumes.PromoCodeResponse{}, fmt.Errorf("failed to decode response: %w", err)
		}

		// Return the data
		return promoCodeResponse, nil

	}
	return attemptRequest()

}

// LoginToBackend logs in to the backend with phoneNumber, login, and password
func LoginToBackend(phoneNumber, login, password string, telegramUserID int64) (volumes.TokenResponse, error) {
	url := getBaseUrl("/api/v1/bot/sign-in")

	// Build the request payload
	payload := volumes.LoginRequest{
		Login:          login,
		Password:       password,
		PhoneNumber:    phoneNumber,
		TelegramUserID: string(telegramUserID),
		UserSession:    helpers.GetUserSession(),
	}

	// Encode payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return volumes.TokenResponse{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create an HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	loggers.Info("response from login to backend", url, " ", req)

	if err != nil {
		return volumes.TokenResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Perform the request
	resp, err := client.Do(req)
	loggers.Info("response from login to backend", err, " ", resp)

	if err != nil {
		return volumes.TokenResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return volumes.TokenResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 response codes
	if resp.StatusCode != http.StatusOK {
		return volumes.TokenResponse{}, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode the JSON response
	var loginResponse volumes.LoginResponse
	if err := json.Unmarshal(body, &loginResponse); err != nil {
		return volumes.TokenResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for a successful response
	if !loginResponse.Success || loginResponse.Status != "OK" {
		return volumes.TokenResponse{}, fmt.Errorf("login failed: server response was not successful")
	}
	// Return the access token
	return loginResponse.Data, nil
}

func GetCategories(language string) ([]volumes.CategoryDataType, error) {
	url := helpers.GetBaseFAQUrl("/api/faqCategory/v1")

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
	var subscriptionResponse volumes.CategoryResponse
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

func GetSubCategories(userData *volumes.UserSession, categoryId, subCategoryId int64) ([]volumes.SubCategoryDataType, error) {
	var apiPath string
	if subCategoryId == -1 {
		apiPath = fmt.Sprintf("/api/faq/v1/withAnswer?categoryId=%d", categoryId)
	} else {
		apiPath = fmt.Sprintf("/api/faq/v1/withAnswer?categoryId=%d&parentFaqId=%d", categoryId, subCategoryId)

	}

	url := helpers.GetBaseFAQUrl(apiPath)
	var attemptRequest func() ([]volumes.SubCategoryDataType, error)
	attemptRequest = func() ([]volumes.SubCategoryDataType, error) {

		// Create HTTP client and request
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			req.Header.Add("Accept", "/")
			return []volumes.SubCategoryDataType{}, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Add("Language", userData.Language)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userData.Token))
		var emptyArray = []volumes.SubCategoryDataType{}
		loggers.Info("response from get subcategories", url, " ", req)

		// Perform the request
		resp, err := client.Do(req)
		loggers.Info("response from subcategories", err, " ", resp)

		if err != nil {
			return emptyArray, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		loggers.Info("response from get user data", resp)
		// If server returns 401, refresh the token and retry
		if resp.StatusCode == http.StatusUnauthorized {
			loggers.Warn("received 401, attempting to refresh token")

			_, refreshErr := RefreshToken(func() volumes.TokenResponse {
				return volumes.TokenResponse{
					AccessToken:  userData.Token,
					TuronId:      userData.TuronId,
					RefreshToken: userData.RefreshToken,
				}
			}, userData)

			if refreshErr != nil {
				return []volumes.SubCategoryDataType{}, fmt.Errorf("failed to refresh token: %w", refreshErr)
			}

			// Retry the request with the new token
			return attemptRequest()
		}

		// Check for non-200 status codes
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body) // Read the body for additional error context
			return emptyArray, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
		}
		// Decode the response
		var subscriptionResponse volumes.SubCategoryResponse
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
	return attemptRequest()
}

func TerminateOwnSession(userData *volumes.UserSession) error {
	apiPath := fmt.Sprintf("/api/v1/users/%d/sessions/terminateOwn", userData.TuronId)
	url := getBaseUrl(apiPath)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Add("Language", userData.Language)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userData.Token))

	loggers.Info("Sending request to terminate session", url, req)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	loggers.Info("Response from terminating session", resp.Body)
	return nil
}

// GetRegions fetches a list of available cities (regions) from the backend.
func GetRegions(userData *volumes.UserSession) ([]volumes.Region, error) {
	url := getBaseUrl("/api/v1/addresses/cities/public?limit=1000&offset=0")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Language", userData.Language)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode JSON response
	var result struct {
		Status  string           `json:"status"`
		Data    []volumes.Region `json:"data"`
		Success bool             `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("unsuccessful response: status = %s", result.Status)
	}

	return result.Data, nil
}

func GetDistricts(userData *volumes.UserSession, cityId int64) ([]volumes.Region, error) {
	// Build URL safely
	base := getBaseUrl("/api/v1/addresses/districts/public")

	// Construct query params properly
	params := fmt.Sprintf("?limit=1000&offset=0&cityId=%d", cityId)
	url := base + params

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Language", userData.Language)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode JSON response
	var result struct {
		Status  string           `json:"status"`
		Data    []volumes.Region `json:"data"`
		Success bool             `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("unsuccessful response: status = %s", result.Status)
	}

	return result.Data, nil
}

func MyApplications(userData *volumes.UserSession, userId int64) ([]volumes.UserApplications, error) {
	// Build URL safely
	base := getBaseUrl(fmt.Sprintf("/api/v1/users/myturonbot/%d/requests", userId))

	// Construct query params properly
	url := base

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Language", userData.Language)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode JSON response
	var result struct {
		Status  string                     `json:"status"`
		Data    []volumes.UserApplications `json:"data"`
		Success bool                       `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("unsuccessful response: status = %s", result.Status)
	}

	return result.Data, nil
}

// SendApplicationToBackend sends a user's application data to the backend
func SendApplicationToBackend(
	regionID int64,
	regionName string,
	districtID int64,
	districtName string,
	fullName string,
	phoneNumber string,
	language string,
	telegramID int64,
	telegramPhoneNumber *string,
) error {

	url := getBaseUrl("/api/v1/users/myturonbot/send-data")

	// 🧩 Build the request payload
	payload := map[string]interface{}{
		"cityId":              regionID,
		"cityName":            regionName,
		"districtId":          districtID,
		"districtName":        districtName,
		"fullName":            fullName,
		"telegramPhoneNumber": phoneNumber,
		"preferredLanguage":   language,
		"telegramId":          telegramID,
	}
	// ✅ Only include optional field if not nil
	if telegramPhoneNumber != nil && *telegramPhoneNumber != "" {
		payload["phoneNumber"] = *telegramPhoneNumber
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 🧩 Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	// 🧩 Decode backend response
	var result struct {
		Status  string `json:"status"`
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("backend returned unsuccessful response: %s", result.Message)
	}
	return nil
}
