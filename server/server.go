package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Simulate server request to get user data
func GetFAQs(credentials string) (*UserData, error) {
	username := credentials[:len(credentials)-2]
	url := fmt.Sprintf("http://example.com/api/user?username=%s", username)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %v", resp.StatusCode)
	}

	var userData UserData
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return nil, err
	}

	return &userData, nil
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
	const url = "http://84.46.247.18/api/v1/internet-tariffs/public?offset=0&limit=100"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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
	fmt.Print("#333333333333333333333333333333333333")
	fmt.Print(err)
	fmt.Print("#333333333333333333333333333333333333")

	if err != nil {
		return nil, err
	}

	return objects.Data, nil
}
