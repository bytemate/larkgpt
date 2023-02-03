package chatgpt

import "testing"

func TestChatGPTRequest(t *testing.T) {
	type args struct {
		msg    string
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				msg:    "Hey",
				userID: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := newChatGPTClient("", "", new(noneMetrics))
			_, err := cli.ChatGPTRequest(tt.args.msg, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatGPTRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
