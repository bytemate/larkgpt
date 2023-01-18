package src

import (
	"github.com/go-resty/resty/v2"
)

// {"response":"Here is an example of CSS animation that changes the background color from red to yellow over a 4-second period, with the first 3 seconds being red and the last 1 second being yellow:\n\n```css\n@keyframes colorChange {\n  0% { background-color: red; }\n  75% { background-color: red; }\n  100% { background-color: yellow; }\n}\n\n#myElement {\n  animation: colorChange 4s;\n}\n```\n\nThis animation is applied to an element with the id \"myElement\", and the animation is called \"colorChange\". The animation lasts for 4 seconds, and the keyframe percentages specify the progression of the animation. At 0% (the start of the animation), the background color is set to red. At 75%, the background color is still red. At 100% (the end of the animation), the background color is set to yellow.\n\nYou can also adjust the timing function for more fluent animation, Like\n```\n#myElement {\n  animation: colorChange 4s ease-in-out;\n}\n```\n\n`ease-in-out` is one of the timing function, this will make the animation start slow and end slow as well.\n","conversationId":"fe46dbb2-9074-4364-8358-abaf08535e5c","messageId":"af3dcb0e-0827-42cc-85fe-b1f9307110f7"}
type ChatGPTResponse struct {
	Response       string `json:"response"`
	ConversationID string `json:"conversationId"`
	MessageID      string `json:"messageId"`
}

func ChatGPTRequest(msg string, userID string) (string, error) {
	client := resty.New()
	if userID != "" {
		client.SetHeader("User-ID", userID)
	}
	resp, err := client.R().
		SetHeaders(
			map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "LarkGPT",
			}).
		SetResult(&ChatGPTResponse{}).
		SetBody(map[string]string{
			"message": msg,
		}).
		Post(Config.ChatGPTAPIURL)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() == 429 {
		return "ChatGPT 访问过于频繁, 请稍后再试.", err
	}
	response := resp.Result().(*ChatGPTResponse)
	if response.Response == "" && resp.Size() == 0 {
		return "", err
	}
	if response != nil {
		return response.Response, nil
	}
	return "", err
}
