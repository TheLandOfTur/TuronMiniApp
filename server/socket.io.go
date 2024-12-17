package server

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/OzodbekX/TuronMiniApp/volumes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	sio "github.com/karagenc/socket.io-go"
	eio "github.com/karagenc/socket.io-go/engine.io"
	"github.com/karagenc/socket.io-go/engine.io/transport"
	"github.com/quic-go/webtransport-go"
	"nhooyr.io/websocket"
)

func ConnectSocketIO(bot *tgbotapi.BotAPI, userSessions *sync.Map, chatID int64) (*sio.ClientSocket, error) {
	// username := pflag.StringP("username", "u", "", "Username")
	// url := pflag.StringP("connect", "c", defaultURL, "URL to connect to")
	// pflag.Parse()

	term, typing, exitFunc, err := initTerm()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	config := &sio.ManagerConfig{
		EIO: eio.ClientConfig{
			RequestHeader: &transport.RequestHeader{},
			UpgradeDone: func(transportName string) {
				fmt.Fprintf(term, "Transport upgraded to: %s\n", transportName)
			},
			HTTPTransport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // DO NOT USE in production. This is for self-signed TLS certificate to work.
				},
			},
			WebSocketDialOptions: &websocket.DialOptions{
				HTTPClient: &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true, // DO NOT USE in production. This is for self-signed TLS certificate to work.
						},
					},
				},
			},
			WebTransportDialer: &webtransport.Dialer{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // DO NOT USE in production. This is for self-signed TLS certificate to work.
				},
			},
		},
	}
	lang := "uz"
	token := ""

	if session, ok := userSessions.Load(chatID); ok {
		user := session.(*volumes.UserSession)
		lang = user.Language
		token = user.Token
	}
	serverURL := os.Getenv("BASE_SOCKET_URL")
	route := serverURL + fmt.Sprintf("/?lang=%s&token=%s", lang, token)
	manager := sio.NewManager(route, config)
	client := manager.Socket("/", nil)
	typing.socket = client

	// if *username == "" {
	// 	fmt.Fprintln(term, "Username cannot be empty. Use -u/--username to set the username.")
	// 	exitFunc(1)
	// }
	if client != nil {
		client.OnConnect(func() {
			fmt.Fprintln(term, "Connected")
		})
		manager.OnError(func(err error) {
			fmt.Fprintf(term, "Error: %v\n", err)
		})
		manager.OnReconnect(func(attempt uint32) {
			fmt.Fprintf(term, "Reconnected. Number of attempts so far: %d\n", attempt)
		})
		client.OnConnectError(func(err any) {
			fmt.Fprintf(term, "Connect error: %v\n", err)
		})
		client.OnEvent("onNewMessageReceived", func(msg volumes.MessageFromSocket) {
			sendOperatorMessage(bot, chatID, msg)
		})
		client.OnEvent("onSuccessSentMessage", func() {
			log.Println("Successfully connected to the onSuccessSentMessage")
		})
		// This will be emitted after the socket is connected.
		client.Emit("add user", "username")

		client.Connect()

		for {
			line, err := term.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(term, "Error: %v\n", err)
				exitFunc(1)
			}
			if line == "" {
				continue
			}
			client.Emit("new message", string(line))
		}
		exitFunc(0)
	}
	return &client, nil
}

func SendMessageToSocketIO(client *sio.ClientSocket, message string, faqId *int64) {
	content := MessageDto{
		Content: message,
		FaqId:   faqId,
	}

	if client != nil {
		(*client).Emit("sendMessageToOperator", content)
	} else {
		log.Println("Socket client is not connected")
	}

}
