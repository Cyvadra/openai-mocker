package main

import (
	"github.com/Cyvadra/openai-mocker/mocker"
)

var sampleRequest = []mocker.ChatRequestMessage{}

func init() {
	sampleRequest = append(sampleRequest, mocker.NewSystemPrompt("You are a helpful agent."))
	sampleRequest = append(sampleRequest, mocker.NewUserMessage("hello world"))
	sampleRequest = append(sampleRequest, mocker.ChatRequestMessage{
		Role:    "assistant",
		Content: "meow",
	})
}

func SampleReply(messages []mocker.ChatRequestMessage) string {
	return messages[len(messages)-1].Content
}

func main() {
	SampleReply(sampleRequest)
	mocker.RunAgent(SampleReply, 3001)
}
