package maxbot

import (
	"context"
	"log"
	"sync"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func NewBot(token string, opt ...Opt) (*Bot, error) {
	ctx, cancel := context.WithCancel(context.Background())
	b := &Bot{
		ctx:      ctx,
		cancel:   cancel,
		onError:  defaultOnError,
		poller:   new(LongPolling),
		handlers: make(map[string]HandlerFunc),
		updates:  make(chan model.Update, 1),
		stop:     make(chan chan struct{}),
	}

	var err error
	for _, o := range opt {
		err = o(b)
		if err != nil {
			return nil, err
		}
	}

	cli, err := maxbot.NewApi(token, b.opts...)
	if err != nil {
		return nil, err
	}

	b.client = cli

	return b, nil
}

type Bot struct {
	ctx    context.Context
	cancel context.CancelFunc
	mx     sync.Mutex
	client *maxbot.Api
	opts   []maxbot.Opt
	poller Poller

	synchronous bool
	onError     func(error, Context)
	handlers    map[string]HandlerFunc
	updates     chan model.Update
	stopBot     chan struct{}
	stop        chan chan struct{}
}

func (b *Bot) Handle(endpoint any, h HandlerFunc, m ...MiddlewareFunc) {
	end := extractEndpoint(endpoint)
	if end == "" {
		panic("maxbot: unsupported endpoint")
	}

	b.handlers[end] = func(c Context) error {
		return applyMiddleware(h, m...)(c)
	}
}

func (b *Bot) Start() {
	if b.poller == nil {
		panic("maxbot: can't start without a poller")
	}

	// do nothing if called twice
	b.mx.Lock()
	if b.stopBot != nil {
		b.mx.Unlock()
		return
	}

	b.stopBot = make(chan struct{})
	b.mx.Unlock()

	stop := make(chan struct{})
	stopConfirm := make(chan struct{})

	go func() {
		b.poller.Poll(b.ctx, b.client.Subscriptions, b.updates)
		close(stopConfirm)
	}()

	for {
		select {
		case upd := <-b.updates:
			b.ProcessUpdate(upd)
		case confirm := <-b.stop:
			close(stop)
			<-stopConfirm
			close(confirm)
			return
		}
	}
}

func (b *Bot) Stop() {
	b.mx.Lock()
	if b.stopBot != nil {
		close(b.stopBot)
		b.stopBot = nil
	}

	b.cancel()
	b.mx.Unlock()

	confirm := make(chan struct{})
	b.stop <- confirm
	<-confirm
}

func (b *Bot) ProcessUpdate(u model.Update) {
	b.ProcessContext(b.NewContext(u))
}

func (b *Bot) NewContext(u model.Update) Context {
	return NewContext(b.client, u)
}

func (b *Bot) ProcessContext(c Context) {
	u := c.Update()

	if maxbot.GetCommand(u) != "" {
		b.handle(maxbot.GetCommand(u), c)
	}

}

func (b *Bot) handle(end string, c Context) bool {
	if handler, ok := b.handlers[end]; ok {
		b.runHandler(handler, c)
		return true
	}
	return false
}

func (b *Bot) runHandler(h HandlerFunc, c Context) {
	f := func() {
		if err := h(c); err != nil {
			b.OnError(err, c)
		}
	}
	if b.synchronous {
		f()
	} else {
		go f()
	}
}

func (b *Bot) OnError(err error, c Context) {
	b.onError(err, c)
}

func extractEndpoint(endpoint interface{}) string {
	switch end := endpoint.(type) {
	case string:
		return end
	case CallbackEndpoint:
		return end.CallbackUnique()
	}
	return ""
}

var defaultOnError = func(err error, c Context) {
	if c != nil {
		log.Println(c.Update().ChatID, err)
	} else {
		log.Println(err)
	}
}
