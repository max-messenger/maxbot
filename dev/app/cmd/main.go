package main

import (
	"context"
	"log"
	"maxbot"
	"os"
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

		err = c.Send(ctx, "msg")
		if err != nil {
			return err
		}

		return nil
	})

	bot.Start()
}
