package server

import (
	"fmt"
	"github.com/OzodbekX/TuronMiniApp/volumes"
	"github.com/zhouhui8915/go-socket.io-client"
	"log"
	"os"
	"sync"
)

func StartSocketIOServer(userSessions *sync.Map, chatID int64) (*socketio_client.Client, error) {
	serverURL := os.Getenv("BASE_SOCKET_URL")
	lang := "uz"
	token := ""

	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		lang = user.Language
		token = user.Token
	}

	opts := &socketio_client.Options{
		Transport: "websocket", // Use WebSocket as the transport
		Query: map[string]string{
			"token": token, // Replace with your token
			"lang":  lang,  // Set your preferred language
		},
	}

	// Connect to the server
	client, err := socketio_client.NewClient(serverURL, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
	}

	if client != nil {
		// Connection success handler
		client.On("connect", func() {
			log.Println("Successfully connected to the server")
		}) // Connection success handler
		client.On("onSuccessSentMessage", func() {
			log.Println("Successfully connected to the onSuccessSentMessage")
		})
		// Connection error handler
		client.On("connect_error", func(err interface{}) {
			log.Printf("Connection error: %v", err)
		})
	}
	fmt.Println("token\n")
	fmt.Println(client)
	fmt.Println("token\n")

	// Keep the client running in the background
	go func() {
		select {} // Keeps the client alive
	}()
	return client, nil // Return the client
}

type MessageDto struct {
	Content string `json:"content"`
	FaqId   *int64 `json:"faqId,omitempty"` // Use a pointer for the optional parameter

}

// Send a message to the server
func SendMessageToServer(client *socketio_client.Client, message string, faqId *int64) {

	content := MessageDto{
		Content: message,
		FaqId:   faqId,
	}

	if client != nil {
		client.Emit("sendMessageToOperator", content)
	} else {
		log.Println("Socket client is not connected")
	}
}
