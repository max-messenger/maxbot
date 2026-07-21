package maxbot

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type HandlerFunc func(Context) error

type Context interface {
	Update() model.Update
	Context() context.Context
	API() *maxbot.Api

	Send(text string, opts ...Option) error
	Answer(text string, opts ...Option) error
	Reply(text string, opts ...Option) error
	Edit(text string, opts ...Option) error
	Delete(opts ...Option) error
}

type nativeContext struct {
	ctx context.Context
	b   *maxbot.Api
	u   model.Update
}

func NewContext(ctx context.Context, b *maxbot.Api, u model.Update) Context {
	return &nativeContext{
		ctx: ctx,
		b:   b,
		u:   u,
	}
}

func (c *nativeContext) Update() model.Update {
	return c.u
}

func (c *nativeContext) Context() context.Context {
	return c.ctx
}

func (c *nativeContext) API() *maxbot.Api {
	return c.b
}

func (c *nativeContext) Send(text string, opts ...Option) error {
	msg := maxbot.NewMessage().
		SetText(text).
		SetUser(c.u.UserID).
		SetChat(c.u.ChatID)

	for _, opt := range opts {
		opt(msg)
	}

	_, err := c.b.Messages.Send(c.ctx, msg)

	return err
}

func (c *nativeContext) Answer(text string, opts ...Option) error {
	msg := maxbot.NewMessage().
		SetText(text).
		SetUser(c.u.UserID).
		SetChat(c.u.ChatID)

	for _, opt := range opts {
		opt(msg)
	}

	mb := msg.MessageBody()
	_, err := c.b.Messages.AnswerOnCallback(c.ctx, c.u.GetCallback().CallbackID, model.CallbackAnswer{Message: &mb})

	return err
}

func (c *nativeContext) Reply(text string, opts ...Option) error {
	msg := maxbot.NewMessage().
		SetUser(c.u.UserID).
		SetChat(c.u.ChatID).
		SetReply(text, c.u.GetMessage().Body.Mid)

	for _, opt := range opts {
		opt(msg)
	}

	_, err := c.b.Messages.Send(c.ctx, msg)

	return err
}

func (c *nativeContext) Edit(text string, opts ...Option) error {
	msg := maxbot.NewMessage().
		SetText(text).
		SetUser(c.u.UserID).
		SetChat(c.u.ChatID)

	for _, opt := range opts {
		opt(msg)
	}

	_, err := c.b.Messages.Send(c.ctx, msg)

	return err
}

func (c *nativeContext) Delete(opts ...Option) error {
	msg := maxbot.NewMessage().
		SetUser(c.u.UserID).
		SetChat(c.u.ChatID)

	for _, opt := range opts {
		opt(msg)
	}

	_, err := c.b.Messages.DeleteMessage(c.ctx, c.u.MessageID)

	return err
}
