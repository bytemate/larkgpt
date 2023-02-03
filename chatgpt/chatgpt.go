package chatgpt

import (
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	// retry
	"github.com/avast/retry-go"
)

type chatGPTClient struct {
	apiHost    string
	apiKey     string
	metricsIns IMetrics
}

func newChatGPTClient(apiHost, apiKey string, metricsIns IMetrics) *chatGPTClient {
	return &chatGPTClient{
		apiHost:    apiHost,
		apiKey:     apiKey,
		metricsIns: metricsIns,
	}
}

func (r *chatGPTClient) ChatGPTRequest(msg string, userID string) (result string, err error) {
	defer func() {
		if err != nil {
			r.metricsIns.EmitChatGPTApiFailed()
		} else {
			r.metricsIns.EmitChatGPTApiSuccess()
		}
	}()

	client := resty.New()
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
		Post(r.apiHost + "message/" + userID)
	if resp.Size() == 0 {
		return "", errors.New("ChatGPT return empty response")
	}
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

func (r *chatGPTClient) ChatGPTOneTimeRequest(msg string) (result string, err error) {
	defer func() {
		if err != nil {
			r.metricsIns.EmitChatGPTApiFailed()
		} else {
			r.metricsIns.EmitChatGPTApiSuccess()
		}
	}()

	client := resty.New()
	var resp *resty.Response
	err = retry.Do(
		func() error {
			var err error
			resp, err = client.R().
				SetHeaders(
					map[string]string{
						"Content-Type": "application/json",
						"User-Agent":   "LarkGPT",
					}).
				SetResult(&ChatGPTResponse{}).
				SetBody(map[string]string{
					"message": msg,
				}).
				Post(r.apiHost + "message")
			if err != nil {
				return err
			}
			if resp.Size() == 0 {
				return errors.New("ChatGPT return empty response")
			}
			return nil
		},
		retry.Attempts(3),
		retry.DelayType(retry.FixedDelay),
		// Delay 3 seconds
		retry.Delay(4*time.Second),
	)
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

func (r *chatGPTClient) DeleteSession(userID string) (err error) {
	defer func() {
		if err != nil {
			r.metricsIns.EmitChatGPTApiFailed()
		} else {
			r.metricsIns.EmitChatGPTApiSuccess()
		}
	}()

	client := resty.New()
	resp, err := client.R().
		SetHeaders(
			map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "LarkGPT",
			}).
		Delete(r.apiHost + "message/" + userID)
	if err != nil {
		return err
	}
	if resp.StatusCode() == 429 {
		return err
	}
	return nil
}

// {"response":"Here is an example of CSS animation that changes the background color from red to yellow over a 4-second period, with the first 3 seconds being red and the last 1 second being yellow:\n\n```css\n@keyframes colorChange {\n  0% { background-color: red; }\n  75% { background-color: red; }\n  100% { background-color: yellow; }\n}\n\n#myElement {\n  animation: colorChange 4s;\n}\n```\n\nThis animation is applied to an element with the id \"myElement\", and the animation is called \"colorChange\". The animation lasts for 4 seconds, and the keyframe percentages specify the progression of the animation. At 0% (the start of the animation), the background color is set to red. At 75%, the background color is still red. At 100% (the end of the animation), the background color is set to yellow.\n\nYou can also adjust the timing function for more fluent animation, Like\n```\n#myElement {\n  animation: colorChange 4s ease-in-out;\n}\n```\n\n`ease-in-out` is one of the timing function, this will make the animation start slow and end slow as well.\n","conversationId":"fe46dbb2-9074-4364-8358-abaf08535e5c","messageId":"af3dcb0e-0827-42cc-85fe-b1f9307110f7"}
type ChatGPTResponse struct {
	Response       string `json:"response"`
	ConversationID string `json:"conversationId"`
	MessageID      string `json:"messageId"`
}
