package listeners

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OzodbekX/TuronMiniApp/handlers/events"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// requestHandler handles incoming HTTP requests to send messages to multiple users.
func RequestHandler(botApi *tgbotapi.BotAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body
		var reqPayload volumes.AlertRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Validate the input
		if len(reqPayload.Messages) == 0 {
			http.Error(w, "No messages provided", http.StatusBadRequest)
			return
		}

		// Iterate over each user-message pair and send messages
		var failedMessages []string
		for _, msg := range reqPayload.Messages {
			if msg.UserID == 0 || msg.Message == "" {
				failedMessages = append(failedMessages, fmt.Sprintf("Invalid data for UserID %d", msg.UserID))
				continue
			}

			// Attempt to send the message
			events.SendMessage(botApi, msg.ChatID, msg.Message)
		}

		// Prepare the response
		if len(failedMessages) > 0 {
			w.WriteHeader(http.StatusPartialContent)
			response := map[string]interface{}{
				"data":    fmt.Sprintf("Some messages failed: %v", failedMessages),
				"status":  "Partial Content",
				"success": false,
			}
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{
				"data":    "All messages sent successfully",
				"status":  "OK",
				"success": true,
			}
			json.NewEncoder(w).Encode(response)
		}
	}
}
