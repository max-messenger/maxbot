package maxbot

const (
	callbackPrefix = "callback_"
	commandPrefix  = "/"
)

const (
	OnMessageCreated        = "message_created"
	OnMessageCallback       = "message_callback"
	OnMessageEdited         = "message_edited"
	OnMessageRemoved        = "message_removed"
	OnBotAdded              = "bot_added"
	OnBotRemoved            = "bot_removed"
	OnUserAdded             = "user_added"
	OnUserRemoved           = "user_removed"
	OnBotStarted            = "bot_started"
	OnBotStopped            = "bot_stopped"
	OnDialogCleared         = "dialog_cleared"
	OnDialogRemoved         = "dialog_removed"
	OnDialogMuted           = "dialog_muted"
	OnDialogUnmuted         = "dialog_unmuted"
	OnChatTitleChangedEvent = "chat_title_changed"
	OnText                  = "text"
)
