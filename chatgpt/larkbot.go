package chatgpt

import (
	"context"
	"log"
	"strings"

	"github.com/chyroc/lark"
)

func isNonsense(msg string) bool {
	// includee message
	// msg include @_all
	return strings.Contains(msg, "@_all") || msg == ""
}

func filterMsg(msg string) string {
	// filter message
	// msg include @_user_1
	if strings.Contains(msg, "@_user_1") {
		msg = strings.ReplaceAll(msg, "@_user_1", "")
	}
	return msg
}

func (r *Client) ReceiveChatGPTMessage(ctx context.Context, msg string, event *lark.EventV2IMMessageReceiveV1) (err error) {
	defer func() {
		if err != nil {
			r.metricsIns.EmitAppFailed()
		} else {
			r.metricsIns.EmitAppSuccess()
		}
	}()

	log.Print("Receive message: ", msg)
	if r.maintained {
		return r.larkIns.replyText(ctx, event.Message.MessageID, "ChatGPT Bot 正在维护中 请稍后重试.请飞书搜索 ChatGPT 讨论群, 选择同款头像进群看进度.")
	}

	var result string
	if event.Message.ChatType == "p2p" {
		result, err = r.chatGPTIns.ChatGPTRequest(msg, event.Sender.SenderID.OpenID)
	} else {
		result, err = r.chatGPTIns.ChatGPTOneTimeRequest(msg)
	}
	log.Println("msg: ", msg, "result: ", result)
	if err != nil {
		log.Println("ChatGPT 请求失败 请稍后重试. ", err)
		return r.larkIns.replyText(ctx, event.Message.MessageID, "ChatGPT 请求失败 请稍后重试.")
	}

	return r.larkIns.replyText(ctx, event.Message.MessageID, result)
}

func (r *Client) ReceiveCommandMessage(ctx context.Context, command string, event *lark.EventV2IMMessageReceiveV1) {
	switch command {
	case "/reset":
		err := r.chatGPTIns.DeleteSession(
			event.Sender.SenderID.OpenID,
		)
		if err != nil {
			r.larkIns.replyText(ctx, event.Message.MessageID, "Reset Failed.")
			return
		}
		r.larkIns.replyText(ctx, event.Message.MessageID, "Reset Success.")
	default:
		r.larkIns.replyText(ctx, event.Message.MessageID, "Unknown Command.")
	}
}

func (r *Client) larkMessageReceiverHandler(ctx context.Context, cli *lark.Lark, schema string, header *lark.EventHeaderV2, event *lark.EventV2IMMessageReceiveV1) (string, error) {
	content, err := lark.UnwrapMessageContent(event.Message.MessageType, event.Message.Content)
	if err != nil {
		return "", err
	}
	msg := ""
	switch event.Message.MessageType {
	case lark.MsgTypeText:
		msg = content.Text.Text
	case lark.MsgTypePost:
		msg = wrapLarkPostMessageText(content)
	default:
		log.Println("暂不支持的消息类型.")
		_ = r.larkIns.replyText(ctx, event.Message.MessageID, "暂不支持的消息类型.")
		return "", nil
	}
	msg = filterMsg(msg)
	if isNonsense(msg) {
		return "", nil
	}
	switch true {
	case strings.HasPrefix(msg, "/"):
		go r.ReceiveCommandMessage(ctx, msg, event)
	default:
		go r.ReceiveChatGPTMessage(ctx, msg, event)
	}
	return "", err
}

func wrapLarkPostMessageText(content *lark.MessageContent) string {
	builder := new(strings.Builder)
	for idx, postContentList := range content.Post.Content {
		if idx != 0 {
			builder.WriteString("\n")
		}
		for _, postContent := range postContentList {
			switch postContent := postContent.(type) {
			case lark.MessageContentPostLink:
				builder.WriteString(postContent.Href)
			case lark.MessageContentPostText:
				builder.WriteString(postContent.Text)
			}
		}
	}
	return builder.String()
}
