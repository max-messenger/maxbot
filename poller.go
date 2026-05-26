package maxbot

import (
	"context"
	"log"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type Client interface {
	GetUpdates(ctx context.Context, marker int64) ([]model.Update, int64, error)
}

type Poller interface {
	Poll(ctx context.Context, c Client, dest chan model.Update)
}

type LongPolling struct {
	Limit        int
	Timeout      time.Duration
	LastUpdateID int
	Marker       int64

	AllowedUpdates []string `yaml:"allowed_updates"`
}

func (p *LongPolling) Poll(ctx context.Context, c Client, dest chan model.Update) {
	var updates []model.Update
	var marker int64
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		updates, marker, err = c.GetUpdates(ctx, marker)
		if err != nil {
			log.Println("GetUpdates: ", err)
			return
		}

		for _, update := range updates {
			dest <- update
		}
	}
}

type MiddlewarePoller struct {
	Capacity int // Default: 1
	Poller   Poller
	//Filter   func(*Update) bool
}

func NewMiddlewarePoller(original Poller, filter func(*model.Update) bool) *MiddlewarePoller {
	return &MiddlewarePoller{
		Poller: original,
		//Filter: filter,
	}
}

//func (p *MiddlewarePoller) Poll(ctx context.Context, c Client, dest chan model.Update) {
//	if p.Capacity < 1 {
//		p.Capacity = 1
//	}
//
//	middle := make(chan model.Update, p.Capacity)
//	stopPoller := make(chan struct{})
//	stopConfirm := make(chan struct{})
//
//	go func() {
//		p.Poller.Poll(b, middle, stopPoller)
//		close(stopConfirm)
//	}()
//
//	for {
//		select {
//		case <-stop:
//			close(stopPoller)
//			<-stopConfirm
//			return
//			//case upd := <-middle:
//			//if p.Filter(&upd) {
//			//	dest <- upd
//			//}
//		}
//	}
//}
