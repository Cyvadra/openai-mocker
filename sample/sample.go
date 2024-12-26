package sample

import "github.com/Cyvadra/openai-mocker/mocker"

func Reply(messages []mocker.ChatRequestMessage) string {
	return messages[len(messages)-1].Content
}
