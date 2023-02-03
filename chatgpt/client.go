package chatgpt

import (
	"fmt"
	"net/http"

	"github.com/chyroc/lark"
)

type Client struct {
	larkIns    *lark.Lark
	chatGPTIns *chatGPTClient
	serverPort int
	maintained bool
}

type ClientConfig struct {
	// Lark
	AppID     string
	AppSecret string

	// ChatGPT
	ChatGPTAPIKey string
	ChatGPTAPIURL string
	Maintained    bool

	// server
	ServerPort int
}

func New(config *ClientConfig) *Client {
	res := new(Client)

	res.larkIns = lark.New(lark.WithAppCredential(config.AppID, config.AppSecret))

	res.chatGPTIns = newChatGPTClient(config.ChatGPTAPIURL, config.ChatGPTAPIKey)

	res.serverPort = config.ServerPort
	if res.serverPort == 0 {
		res.serverPort = 9726
	}

	res.maintained = config.Maintained

	return res
}

func (r *Client) Start() error {
	r.larkIns.EventCallback.HandlerEventV2IMMessageReceiveV1(r.ReciverMessage)

	http.HandleFunc("/event", func(w http.ResponseWriter, req *http.Request) {
		r.larkIns.EventCallback.ListenCallback(req.Context(), req.Body, w)
	})

	fmt.Printf("start server: %d ...\n", r.serverPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", r.serverPort), nil)
}
