package tg_client

import (
	"testing"
)

func Test_buildURL(t *testing.T) {
	type args struct {
		tgClient *tgClient
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "base test",
			args: args{
				tgClient: &tgClient{
					host:     "api.telegram.org",
					basePath: "bot1test_token",
				},
			},
			want: "https://api.telegram.org/bot1test_token/sendMessage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildURL(tt.args.tgClient); got != tt.want {
				t.Errorf("buildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
