package maxbot

type CallbackEndpoint interface {
	CallbackUnique() string
}
