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

		// the bot will be subscribed to the specified types
		//maxbot.WithWebhook("http://my-bot.cloud.hooli.local/webhook", "secret", []string{
		//		maxbot.OnBotAdded,
		//		maxbot.OnMessageCreated,
		//		maxbot.OnMessageCallback,
		//	}),
	}

	token := os.Getenv("BOT_TOKEN")

	bot, err := maxbot.NewApi(token, opts...)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/help", func(c maxbot.Context) error {
		kb := model.NewKeyboard()
		kb.AddRow().
			AddLink("Документация", "https://dev.max.ru/docs").
			AddCallBack("нажми на меня", "pushBtn")

		return c.Send("Основная информация:", maxbot.WithKeyboard(kb))
	})

	bot.Handle("/command", func(c maxbot.Context) error {
		command := c.Update().GetCommand()
		msg := fmt.Sprintf(
			"command: %s\nbot name: %s\n params: \n%s\n text: %s\n",
			command.Command, command.BotName,
			strings.Join(command.Params, "\n"),
			command.RemainingText,
		)

		return c.Send(msg)
	})

	bot.Handle("/reply", func(c maxbot.Context) error {
		kb := model.NewKeyboard()
		kb.AddRow().
			AddLink("Документация", "https://dev.max.ru/docs")

		return c.Reply("reply", maxbot.WithKeyboard(kb))
	})

	bot.HandleCallback("pushBtn", func(c maxbot.Context) error {
		kb := model.NewKeyboard()
		kb.AddRow().
			AddLink("Документация", "https://dev.max.ru/docs")

		// c.Answer("Изменено") // не только поменяет текст, но и удалит все кнопки
		return c.Answer("Изменено", maxbot.WithKeyboard(kb))
	})

	//bot.HandleCallback("pushBtn", func(c maxbot.Context) error {
	//	return c.Edit("Изменено") // изменить текст + удалить кнопки
	//	return c.Edit("Изменено", maxbot.WithAttachments(c.Update().GetMessage().Body.Attachments)) // изменить только текст
	//})

	bot.Handle(maxbot.OnChatTitleChangedEvent, func(c maxbot.Context) error {
		return c.Send("Заголовок чата изменен")
	})

	bot.Handle(maxbot.OnMessageCreated, func(c maxbot.Context) error {
		//err = c.Send(fmt.Sprintf("%s - принято", c.Update().GetMessage().Body.Text))
		//if err != nil {
		//	return err
		//}

		fmt.Println("-->", c.Update().GetMessage().Body.Text, c.Update().MessageID)

		return nil
	})

	bot.Handle(maxbot.OnCommentCreated, func(c maxbot.Context) error {
		fmt.Printf(
			"new comment on post_id: %s, comment_id: %s, comment: %s\n",
			c.Update().PostID,
			c.Update().CommentID,
			c.Update().Message.Body.Text,
		)
		return nil
	})

	bot.Start()
}
