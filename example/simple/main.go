package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
	"github.com/max-messenger/maxbot"
)

func main() {
	opts := []maxbot.Opt{
		maxbot.WithHTTPClient(&http.Client{Timeout: 25 * time.Second}),
	}

	token := os.Getenv("BOT_TOKEN")

	bot, err := maxbot.NewApi(token, opts...)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/info", func(c maxbot.Context) error {
		kb := model.NewKeyboard()
		kb.AddRow().AddLink("docs", "https://dev.max.ru/docs")
		err = c.Send("max мне в руки", maxbot.WithKeyboard(kb))
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle(maxbot.OnChatTitleChangedEvent, func(c maxbot.Context) error {
		err = c.Send("title changed")
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle(maxbot.OnText, func(c maxbot.Context) error {
		err = c.Send(fmt.Sprintf("%s - принято", c.Update().GetMessage().Body.Text))
		if err != nil {
			return err
		}

		return nil
	})

	bot.Start()
}
