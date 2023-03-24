package larkgpt

import (
	"context"

	"github.com/chyroc/lark"
	"github.com/chyroc/lark/card"
)

func (r *Client) replyMaintained(msgID string) error {
	if !r.enableCardResp {
		text := "ChatGPT Bot 正在维护中 请稍后重试.请飞书搜索 ChatGPT 讨论群, 选择同款头像进群看进度."
		return r.larkIns.replyText(context.Background(), msgID, text)
	}
	cardIns := card.Card(
		card.Div().SetFields(
			card.FieldMarkdown("ChatGPT Bot 正在维护中 请稍后重试."),
		),
		card.HR(),
		card.Note(card.Text("请飞书搜索 ChatGPT 讨论群, 选择同款头像进群看进度.")),
	).SetHeader(
		card.Header("ChatGPT Bot Error").SetRed(),
	)
	return r.larkIns.replyCard(context.Background(), msgID, cardIns)
}

func (r *Client) replyChatGPTError(msgID string, text string) error {
	if !r.enableCardResp {
		return r.larkIns.replyText(context.Background(), msgID, text)
	}
	cardIns := card.Card(
		card.Div().SetFields(
			card.FieldMarkdown(text),
		),
	).SetHeader(
		card.Header("ChatGPT Bot Error").SetRed(),
	)
	return r.larkIns.replyCard(context.Background(), msgID, cardIns)
}

func (r *Client) replyChatGPTSuccess(msgID string, text string) error {
	if !r.enableCardResp {
		return r.larkIns.replyText(context.Background(), msgID, text)
	}
	cardIns := card.Card(
		card.Div().SetFields(
			card.FieldMarkdown(text),
		),
	).SetHeader(
		card.Header("ChatGPT Bot Error").SetGreen(),
	)
	return r.larkIns.replyCard(context.Background(), msgID, cardIns)
}

func (r *larkClient) replyChatGPTMessage(msgID, title, text string, enableCardResp bool) error {
	if !enableCardResp {
		return r.replyText(context.Background(), msgID, text)
	}
	cardIns := card.Card(
		card.Div().SetFields(
			card.FieldMarkdown(text),
		),
		card.HR(),
		card.Note(
			card.I18nText(&lark.I18NText{
				ZhCn: "HELP: 回复消息支持上下文对话; /reset 重置对话",
				EnUs: "HELP: Reply messages with context; send /reset resets dialogue",
			}),
		),
	).SetHeader(card.Header(title))
	return r.replyCard(context.Background(), msgID, cardIns)
}
