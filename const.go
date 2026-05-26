package maxbot

import "github.com/max-messenger/max-bot-api-client-go/v2/model"

const (
	OnMessageCreated        = string(model.UpdateMessageCreated)
	OnMessageCallback       = string(model.UpdateMessageCallback)
	OnMessageEdited         = string(model.UpdateMessageEdited)
	OnMessageRemoved        = string(model.UpdateMessageRemoved)
	OnBotAdded              = string(model.UpdateBotAdded)
	OnBotRemoved            = string(model.UpdateBotRemoved)
	OnUserAdded             = string(model.UpdateUserAdded)
	OnUserRemoved           = string(model.UpdateUserRemoved)
	OnBotStarted            = string(model.UpdateBotStarted)
	OnBotStopped            = string(model.UpdateBotStopped)
	OnDialogCleared         = string(model.UpdateDialogCleared)
	OnDialogRemoved         = string(model.UpdateDialogRemoved)
	OnDialogMuted           = string(model.UpdateDialogMuted)
	OnDialogUnmuted         = string(model.UpdateDialogUnmuted)
	OnChatTitleChangedEvent = string(model.UpdateChatTitleChanged)
)
