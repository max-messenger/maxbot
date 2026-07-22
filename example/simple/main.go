package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

	bot.Handle("/help", func(c maxbot.Context) error {
		kb := model.NewKeyboard()
		kb.AddRow().
			AddLink("docs", "https://dev.max.ru/docs").
			AddCallBack("push me baby", "pushBtn")

		err = c.Send("max мне в руки", maxbot.WithKeyboard(kb))
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle("/command", func(c maxbot.Context) error {
		command := c.Update().GetCommand()
		msg := fmt.Sprintf(
			"command: %s\nbot name: %s\n params: \n%s\n text: %s\n",
			command.Command, command.BotName,
			strings.Join(command.Params, "\n"),
			command.RemainingText,
		)
		err = c.Send(msg)
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle("/reply", func(c maxbot.Context) error {
		kb := model.NewKeyboard()
		kb.AddRow().
			AddLink("docs", "https://dev.max.ru/docs")

		err = c.Reply("reply", maxbot.WithKeyboard(kb))
		if err != nil {
			return err
		}

		return nil
	})

	bot.HandleCallback("pushBtn", func(c maxbot.Context) error {
		kb := model.NewKeyboard()
		kb.AddRow().
			AddLink("docs", "https://dev.max.ru/docs")

		err = c.Answer("surprise", maxbot.WithKeyboard(kb))
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
		//err = c.Send(fmt.Sprintf("%s - принято", c.Update().GetMessage().Body.Text))
		//if err != nil {
		//	return err
		//}
		fmt.Println("-->", c.Update().GetMessage().Body.Text)

		return nil
	})

	bot.Start()
}
