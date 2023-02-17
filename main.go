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
	godotenv.Load(".env", "../.env")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &larkgpt.ClientConfig{
		AppID:                     os.Getenv("APP_ID"),
		AppSecret:                 os.Getenv("APP_SECRET"),
		LarkOpenBaseURL:           os.Getenv("LARK_OPEN_BASE_URL"), // default is https://open.feishu.cn, for larksuite, please use https://open.larksuite.com
		LarkWWWBaseURL:            os.Getenv("LARK_WWW_BASE_URL"),  // default is https://www.feishu.cn, for larksuite, please use https://www.larksuite.com
		ChatGPTAPIKey:             os.Getenv("CHATGPT_API_KEY"),
		ChatGPTAPIURL:             os.Getenv("CHATGPT_API_URL"),
		ServerPort:                port,
		Maintained:                os.Getenv("MAINTAINED") == "true",
		EnableSessionForLarkGroup: os.Getenv("ENABLE_SESSION_FOR_LARK_GROUP") == "true",
	}, nil
}
