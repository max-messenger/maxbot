package maxbot

import (
	"regexp"
)

var reCommand = regexp.MustCompile(`^/([a-zA-Z_]+)`)

func GetCommand(s string) string {
	match := reCommand.FindAllString(s, -1)
	if len(match) > 0 {
		return match[0]
	}

	return ""
}

var reCommandCallback = regexp.MustCompile(`^([a-zA-Z_]+)`)

func GetCallbackCommand(s string) string {
	match := reCommandCallback.FindAllString(s, -1)
	if len(match) > 0 {
		return match[0]
	}

	return ""
}
