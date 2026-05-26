package maxbot

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type HandlerFunc func(Context) error

type Context interface {
	Update() model.Update
	Send(ctx context.Context, text string, opts ...Option) error
}

type nativeContext struct {
	b *maxbot.Api
	u model.Update
}

func NewContext(b *maxbot.Api, u model.Update) Context {
	return &nativeContext{

		b: b,
		u: u,
	}
}

func (c *nativeContext) Update() model.Update {
	return c.u
}

func (c *nativeContext) Send(ctx context.Context, text string, opts ...Option) error {
	msg := maxbot.NewMessage().
		SetText(text).
		SetUser(c.u.UserID).
		SetChat(c.u.ChatID)

	for _, opt := range opts {
		opt(msg)
	}

	_, err := c.b.Messages.Send(ctx, msg)

	return err
}
