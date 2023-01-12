package src

import (
	"context"
	"log"

	"github.com/chyroc/lark"
)

var LarkServer *lark.Lark

func init() {
	LarkServer = lark.New(lark.WithAppCredential(Config.AppID, Config.AppSecret))
	LarkServer.EventCallback.HandlerEventV2IMMessageReceiveV1(ReciverMessage)
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
	go func() {
		log.Print("Receive message: ", msg)
		result, err := ChatGPTRequest(msg, event.Sender.SenderID.OpenID)
		log.Println("msg: ", msg, "result: ", result)
		if err != nil {
			_, _, err = cli.Message.Reply(event.Message.MessageID).SendText(context.Background(), "ChatGPT Bot 调用失败 请稍后重试. "+"error: "+err.Error())
			if err != nil {
				log.Println("LarkAPI 调用失败 请稍后重试. ", err)
			}
		} else {
			_, _, err = cli.Message.Reply(event.Message.MessageID).SendText(context.Background(), result)
			if err != nil {
				log.Println("LarkAPI 调用失败 请稍后重试. ", err)
			}
		}
	}()
	return "", err
}
