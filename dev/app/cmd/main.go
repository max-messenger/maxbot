package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
	"github.com/max-messenger/maxbot"
)

func main() {
	ctx := context.Background()
	opts := []maxbot.Opt{
		maxbot.WithHTTPClient(newHttpClient()),
	}

	token := os.Getenv("BOT_TOKEN")

	bot, err := maxbot.NewBot(token, opts...)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/info", func(c maxbot.Context) error {
		kb := model.NewKeyboard()
		kb.AddRow().AddLink("ya", "https://ya.ru")
		err = c.Send(ctx, "fx мне в руки", maxbot.WithKeyboard(kb))
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle(maxbot.OnChatTitleChangedEvent, func(c maxbot.Context) error {
		err = c.Send(ctx, "title changed")
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle(maxbot.OnText, func(c maxbot.Context) error {
		err = c.Send(ctx, fmt.Sprintf("%s - сам такой", c.Update().GetMessage().Body.Text))
		if err != nil {
			return err
		}

		return nil
	})

	bot.Start()
}
