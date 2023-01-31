package src

// Go load .env file
import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	// Lark
	AppID     string
	AppSecret string
	// ChatGPT
	ChatGPTAPIKey string
	ChatGPTAPIURL string
	Maintained    bool
}

var Config AppConfig

func init() {
	_ = godotenv.Load(
		".env",
		"../.env",
	)
	Config = AppConfig{
		AppID:         os.Getenv("APP_ID"),
		AppSecret:     os.Getenv("APP_SECRET"),
		ChatGPTAPIKey: os.Getenv("CHATGPT_API_KEY"),
		ChatGPTAPIURL: os.Getenv("CHATGPT_API_URL"),
		Maintained:    os.Getenv("MAINTAINED") == "true",
	}
}
