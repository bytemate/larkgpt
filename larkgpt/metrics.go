package larkgpt

type IMetrics interface {
	// larkgpt api
	EmitChatGPTApiFailed()
	EmitChatGPTApiSuccess()

	// lark api
	EmitLarkApiFailed()
	EmitLarkApiSuccess()

	// app
	EmitAppSuccess()
	EmitAppFailed()
}

type noneMetrics struct{}

func (r *noneMetrics) EmitChatGPTApiFailed() {}

func (r *noneMetrics) EmitChatGPTApiSuccess() {}

func (r *noneMetrics) EmitLarkApiFailed() {}

func (r *noneMetrics) EmitLarkApiSuccess() {}

func (r *noneMetrics) EmitAppSuccess() {}

func (r *noneMetrics) EmitAppFailed() {}
