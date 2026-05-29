package maxbot

import (
	"context"
	"fmt"
	"testing"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestApi_ProcessContext(t *testing.T) {

	type fields struct {
		ctx         context.Context
		cancel      context.CancelFunc
		client      *maxbot.Api
		conf        conf
		group       *Group
		synchronous bool
		onError     func(error, Context)
		handlers    map[string]HandlerFunc
		updates     chan model.Update
		Info        model.BotInfo
	}
	type args struct {
		c Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"callback test:1",
			fields{
				handlers: map[string]HandlerFunc{
					"callback_test": func(Context) error {
						fmt.Println("callback test:1")
						return nil
					},
				},
			},
			args{NewContext(context.Background(), nil, model.Update{
				Callback: &model.Callback{
					Payload: "test:1",
				},
				Message: &model.MessageUpdate{},
			})},
		},
		{
			"callback json test1:{id:1}",
			fields{
				handlers: map[string]HandlerFunc{
					"callback_test": func(Context) error {
						fmt.Println("callback json test1:{id:1}")
						return nil
					},
				},
			},
			args{NewContext(context.Background(), nil, model.Update{
				Callback: &model.Callback{
					Payload: "test:{id:1}",
				},
				Message: &model.MessageUpdate{},
			})},
		},
		{
			"command /test",
			fields{
				handlers: map[string]HandlerFunc{
					"/test": func(Context) error {
						fmt.Println("command /test")
						return nil
					},
				},
			},
			args{NewContext(context.Background(), nil, model.Update{
				Message: &model.MessageUpdate{
					Body: model.MessageBody{
						Text: "/test",
					},
				},
			})},
		},
		{
			"command /test:1",
			fields{
				handlers: map[string]HandlerFunc{
					"/test": func(Context) error {
						fmt.Println("command /test:1")
						return nil
					},
				},
			},
			args{NewContext(context.Background(), nil, model.Update{
				Message: &model.MessageUpdate{
					Body: model.MessageBody{
						Text: "/test:1",
					},
				},
			})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Api{
				ctx:         tt.fields.ctx,
				cancel:      tt.fields.cancel,
				client:      tt.fields.client,
				conf:        tt.fields.conf,
				group:       tt.fields.group,
				synchronous: tt.fields.synchronous,
				onError:     tt.fields.onError,
				handlers:    tt.fields.handlers,
				updates:     tt.fields.updates,
				Info:        tt.fields.Info,
			}
			a.ProcessContext(tt.args.c)
		})
	}
}
