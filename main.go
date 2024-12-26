package main

import (
	"github.com/Cyvadra/openai-mocker/mocker"
)

func main() {
	mocker.RunAgent(mocker.SampleReply, 3001)
}
