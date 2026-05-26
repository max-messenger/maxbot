package maxbot

import (
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
)

type Opt func(cli *Bot) error

func WithHTTPClient(cli maxbot.HttpClient) Opt {
	return func(b *Bot) error {
		b.opts = append(b.opts, maxbot.WithHTTPClient(cli))

		return nil
	}
}

func WithBaseURL(baseURL string) Opt {
	return func(b *Bot) error {
		b.opts = append(b.opts, maxbot.WithBaseURL(baseURL))

		return nil
	}
}

func WithPollingPause(d time.Duration) Opt {
	return func(b *Bot) error {
		b.opts = append(b.opts, maxbot.WithPollingPause(d))

		return nil
	}
}

func WithPollingTimeout(d time.Duration) Opt {
	return func(b *Bot) error {
		b.opts = append(b.opts, maxbot.WithPollingTimeout(d))

		return nil
	}
}

func WithPoller(poller Poller) Opt {
	return func(b *Bot) error {
		b.poller = poller

		return nil
	}
}
