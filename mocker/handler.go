package mocker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const responseId string = "chatcmpl-openai_mocker"
const responseObject string = "chat.completion.chunk"

var responseCompletionId string = "chatcmpl-7f8Qxn9XkoGsVcl0RVGltZpPeqMAG"

func RunStreamAgentOnPath(customHandler StreamHandler, port int, path string) {
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	stopReason := "stop"
	server := gin.Default()
	server.Use(CORS())
	server.POST(path, func(c *gin.Context) {
		var chatRequest ChatRequest
		if err := c.ShouldBindJSON(&chatRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if chatRequest.Stream {
			SetEventStreamHeaders(c)
			streamResponse := ChatCompletionsStreamResponse{
				Id:      responseId,
				Object:  responseObject,
				Created: time.Now().Unix(),
				Model:   chatRequest.Model,
			}
			streamResponseChoice := ChatCompletionsStreamResponseChoice{}
			c.Stream(func(w io.Writer) bool {
				dataChan, stopChain := customHandler(chatRequest.Messages)
				if dataChan == nil {
					return false
				}
				select {
				case s := <-dataChan:
					streamResponseChoice.Delta.Content = string(s)
					streamResponse.Choices = []ChatCompletionsStreamResponseChoice{streamResponseChoice}
					jsonStr, _ := json.Marshal(streamResponse)
					c.Render(-1, CustomEvent{Data: "data: " + string(jsonStr)})
					return true
				case <-stopChain:
					streamResponseChoice.FinishReason = &stopReason
					c.Render(-1, CustomEvent{Data: "data: [DONE]"})
					return false
				}
			})
		} else {
			response := Stream2String(customHandler(chatRequest.Messages))
			c.JSON(http.StatusOK, Completion{
				Id:      responseCompletionId,
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   chatRequest.Model,
				Choices: []Choice{
					{
						Index: 0,
						Message: Message{
							Role:    "assistant",
							Content: response,
						},
						FinishReason: "length",
					},
				},
				Usage: Usage{
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

func RunAgentOnPath(customHandler Handler, port int, path string) {
	RunStreamAgentOnPath(func(messages []ChatRequestMessage) (dataChan chan string, stopChan chan bool) {
		return String2Stream(customHandler(messages))
	}, port, path)
}

func RunAgent(customHandler Handler, port int) {
	RunAgentOnPath(customHandler, port, "/v1/chat/completions")
}

func RunStreamAgent(customHandler StreamHandler, port int) {
	RunStreamAgentOnPath(customHandler, port, "/v1/chat/completions")
}
