package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Cyvadra/openai-mocker/mocker"
	"github.com/gin-gonic/gin"
)

var port int = 3001

func main() {
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	stopReason := "stop"
	server := gin.Default()
	server.Use(mocker.CORS())
	server.POST("/v1/chat/completions", func(c *gin.Context) {
		var chatRequest mocker.ChatRequest
		if err := c.ShouldBindJSON(&chatRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		prompt := "This is a mock server."
		if len(chatRequest.Messages) != 0 {
			prompt = chatRequest.Messages[len(chatRequest.Messages)-1].Content
		}
		response := mocker.Prompt2response(prompt)

		if chatRequest.Stream {
			mocker.SetEventStreamHeaders(c)
			dataChan := make(chan string)
			stopChan := make(chan bool)
			streamResponse := mocker.ChatCompletionsStreamResponse{
				Id:      "chatcmpl-openai_mocker",
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   "gpt-3.5-turbo",
			}
			streamResponseChoice := mocker.ChatCompletionsStreamResponseChoice{}
			go func() {
				for i, s := range response {
					streamResponseChoice.Delta.Content = string(s)
					if i == len(response)-1 {
						streamResponseChoice.FinishReason = &stopReason
					}
					streamResponse.Choices = []mocker.ChatCompletionsStreamResponseChoice{streamResponseChoice}
					jsonStr, _ := json.Marshal(streamResponse)
					dataChan <- string(jsonStr)
				}
				stopChan <- true
			}()

			c.Stream(func(w io.Writer) bool {
				select {
				case data := <-dataChan:
					c.Render(-1, mocker.CustomEvent{Data: "data: " + data})
					return true
				case <-stopChan:
					c.Render(-1, mocker.CustomEvent{Data: "data: [DONE]"})
					return false
				}
			})
		} else {
			c.JSON(http.StatusOK, mocker.Completion{
				Id:      "chatcmpl-7f8Qxn9XkoGsVcl0RVGltZpPeqMAG",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "gpt-3.5-turbo",
				Choices: []mocker.Choice{
					{
						Index: 0,
						Message: mocker.Message{
							Role:    "assistant",
							Content: prompt,
						},
						FinishReason: "length",
					},
				},
				Usage: mocker.Usage{
					PromptTokens:     9,
					CompletionTokens: 1,
					TotalTokens:      10,
				},
			})
		}
	})

	log.Printf("Starting server on port %d", port)
	log.Fatal(server.Run(":" + strconv.Itoa(port)))
}
