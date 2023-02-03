package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/bytemate/larkgpt/larkgpt"
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
	err := godotenv.Load(".env", "../.env")
	if err != nil {
		return nil, err
	}
	return &larkgpt.ClientConfig{
		AppID:         os.Getenv("APP_ID"),
		AppSecret:     os.Getenv("APP_SECRET"),
		ChatGPTAPIKey: os.Getenv("CHATGPT_API_KEY"),
		ChatGPTAPIURL: os.Getenv("CHATGPT_API_URL"),
		Maintained:    os.Getenv("MAINTAINED") == "true",
	}, nil
}
