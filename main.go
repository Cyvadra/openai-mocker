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
	"github.com/Cyvadra/openai-mocker/sample"
	"github.com/gin-gonic/gin"
)

const port int = 3001
const responseId string = "chatcmpl-openai_mocker"
const responseObject string = "chat.completion.chunk"

var responseCompletionId string = "chatcmpl-7f8Qxn9XkoGsVcl0RVGltZpPeqMAG"

func main() {
	RunAgent(sample.Reply, port)
}

func RunAgent(customHandler mocker.Handler, port int) {
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

		// generate response
		response := customHandler(chatRequest.Messages)

		if chatRequest.Stream {
			mocker.SetEventStreamHeaders(c)
			dataChan := make(chan string)
			stopChan := make(chan bool)
			streamResponse := mocker.ChatCompletionsStreamResponse{
				Id:      responseId,
				Object:  responseObject,
				Created: time.Now().Unix(),
				Model:   chatRequest.Model,
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
				Id:      responseCompletionId,
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   chatRequest.Model,
				Choices: []mocker.Choice{
					{
						Index: 0,
						Message: mocker.Message{
							Role:    "assistant",
							Content: "Non-stream reply not implemented!",
						},
						FinishReason: "length",
					},
				},
				Usage: mocker.Usage{
					PromptTokens:     0,
					CompletionTokens: 0,
					TotalTokens:      0,
				},
			})
		}
	})

	log.Printf("Starting server on port %d", port)
	log.Fatal(server.Run(":" + strconv.Itoa(port)))
}
