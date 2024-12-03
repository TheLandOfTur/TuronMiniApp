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
func getUserDataFromServer(credentials string) (*UserData, error) {
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
