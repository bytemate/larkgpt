package main

import (
	"log"
	"os"

	"github.com/bytemate/larkgpt/larkgpt"
	"github.com/joho/godotenv"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := larkgpt.New(config)

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() (*larkgpt.ClientConfig, error) {
	godotenv.Load(".env", "../.env")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &larkgpt.ClientConfig{
		AppID:                     os.Getenv("APP_ID"),
		AppSecret:                 os.Getenv("APP_SECRET"),
		ChatGPTAPIKey:             os.Getenv("CHATGPT_API_KEY"),
		ChatGPTAPIURL:             os.Getenv("CHATGPT_API_URL"),
		ServerPort:                port,
		Maintained:                os.Getenv("MAINTAINED") == "true",
		EnableSessionForLarkGroup: os.Getenv("ENABLE_SESSION_FOR_LARK_GROUP") == "true",
		EnableCardResp:            os.Getenv("ENABLE_CARD_RESP") == "true",
	}, nil
}
