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

func WithFormat(format model.TextFormat) Option {
	return func(msg *maxbot.Message) {
		msg.SetFormat(format)
	}
}

func WithChat(id int64) Option {
	return func(msg *maxbot.Message) {
		msg.SetChat(id)
	}
}

func WithUser(id int64) Option {
	return func(msg *maxbot.Message) {
		msg.SetUser(id)
	}
}

// WithAttachByToken создает опцию для добавления вложения к сообщению по токену файла.
//
// Используется паттерн функциональных опций (functional options pattern) для гибкого
// конструирования сообщений с вложениями различных типов (фото, видео, документы и т.д.).
//
// Параметры:
//   - fileToken: строковый токен файла (например, URL или идентификатор загруженного файла)
//   - at: тип вложения из модели AttachmentType (Photo, Video, Document и др.)
//
// Возвращает:
//
//	Option — функцию-модификатор, которая добавляет вложение к сообщению.
//
// Пример использования:
//
//		// Добавление фото к сообщению
//	 msg := maxbot.NewMessage()
//	 msg.WithAttachByToken("photo_token_123",  model.AttachImage)
//
//		// Или через цепочку опций при построении сообщения
//		message := BuildMessage(
//		    WithText("Привет!"),
//		    WithAttachByToken("file_token", model.AttachImage),
//		)
func WithAttachByToken(fileToken string, at model.AttachmentType) Option {
	return func(msg *maxbot.Message) {
		msg.AddAttachByToken(fileToken, at)
	}
}
