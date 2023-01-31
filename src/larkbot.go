package src

import (
	"context"
	"log"
	"strings"

	"github.com/chyroc/lark"
)

var LarkServer *lark.Lark

func init() {
	LarkServer = lark.New(lark.WithAppCredential(Config.AppID, Config.AppSecret))
	LarkServer.EventCallback.HandlerEventV2IMMessageReceiveV1(ReciverMessage)
}
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
func ReciverChatGPTMessage(msg string, cli *lark.Lark, event *lark.EventV2IMMessageReceiveV1) error {
	log.Print("Receive message: ", msg)
	if Config.Maintained {
		_, _, err := cli.Message.Reply(event.Message.MessageID).SendText(context.Background(), "ChatGPT Bot 正在维护中 请稍后重试.请飞书搜索 ChatGPT 讨论群, 选择同款头像进群看进度.")
		if err != nil {
			log.Println("LarkAPI 调用失败 请稍后重试. ", err)
		}
		return nil
	}
	var result string
	var err error
	if event.Message.ChatType == "p2p" {
		result, err = ChatGPTRequest(msg, event.Sender.SenderID.OpenID)
	} else {
		result, err = ChatGPTOneTimeRequest(msg)
	}
	log.Println("msg: ", msg, "result: ", result)
	if err != nil {
		log.Println("ChatGPT 请求失败 请稍后重试. ", err)
		_, _, err := cli.Message.Reply(event.Message.MessageID).SendText(context.Background(), "ChatGPT 请求失败 请稍后重试.")
		if err != nil {
			log.Println("LarkAPI 调用失败 请稍后重试. ", err)
		}
		return nil
	}
	_, _, err = cli.Message.Reply(event.Message.MessageID).SendText(context.Background(), result)
	if err != nil {
		log.Println("LarkAPI 调用失败 请稍后重试. ", err)
	}
	return nil
}

func ReciverMessage(ctx context.Context, cli *lark.Lark, schema string, header *lark.EventHeaderV2, event *lark.EventV2IMMessageReceiveV1) (string, error) {
	content, err := lark.UnwrapMessageContent(event.Message.MessageType, event.Message.Content)
	if err != nil {
		return "", err
	}
	msg := ""
	switch event.Message.MessageType {
	case lark.MsgTypeText:
		msg = content.Text.Text
	default:
		log.Println("暂不支持的消息类型.")
		_, _, _ = cli.Message.Reply(event.Message.MessageID).SendText(ctx, "暂不支持的消息类型.")
		return "", nil
	}
	msg = filterMsg(msg)
	if isNonsense(msg) {
		return "", nil
	}
	go ReciverChatGPTMessage(msg, cli, event)
	return "", err
}
