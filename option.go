package maxbot

import (
	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type Option func(msg *maxbot.Message)

func WithKeyboard(keyboard *model.Keyboard) Option {
	return func(msg *maxbot.Message) {
		msg.AddKeyboard(keyboard)
	}
}
