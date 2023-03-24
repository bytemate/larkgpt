package larkgpt

import (
	"context"
	"log"

	"github.com/chyroc/lark"
)

type larkClient struct {
	cli        *lark.Lark
	metricsIns IMetrics
}

func newLarkClient(cli *lark.Lark, metricsIns IMetrics) *larkClient {
	return &larkClient{
		cli:        cli,
		metricsIns: metricsIns,
	}
}

func (r *larkClient) replyText(ctx context.Context, msgID, text string) error {
	_, _, err := r.cli.Message.Reply(msgID).SendText(ctx, text)
	if err != nil {
		r.metricsIns.EmitLarkApiFailed()
		log.Println("LarkAPI 调用失败 请稍后重试. ", err)
	} else {
		r.metricsIns.EmitLarkApiSuccess()
	}
	return err
}

func (r *larkClient) replyCard(ctx context.Context, msgID string, card *lark.MessageContentCard) error {
	_, _, err := r.cli.Message.Reply(msgID).SendCard(ctx, card.String())
	if err != nil {
		r.metricsIns.EmitLarkApiFailed()
		log.Println("LarkAPI 调用失败 请稍后重试. ", err)
	} else {
		r.metricsIns.EmitLarkApiSuccess()
	}
	return err
}
