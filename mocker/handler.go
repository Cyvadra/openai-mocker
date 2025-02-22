package mocker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const responseId string = "chatcmpl-openai_mocker"
const responseObject string = "chat.completion.chunk"

var responseCompletionId string = "chatcmpl-7f8Qxn9XkoGsVcl0RVGltZpPeqMAG"

func GenerateStreamAgent(customHandler StreamHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		var chatRequest ChatRequest
		stopReason := "stop"
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
						FinishReason: "stop",
					},
				},
				Usage: Usage{
					PromptTokens:     1,
					CompletionTokens: len(response),
					TotalTokens:      1 + len(response),
				},
			})
		}
	}
}

func RegisterStreamAgentOnPath(customHandler StreamHandler, server *gin.Engine, path string) {
	server.POST(path, GenerateStreamAgent(customHandler))
}

func RegisterAgentOnPath(customHandler Handler, server *gin.Engine, path string) {
	RegisterStreamAgentOnPath(func(messages []ChatRequestMessage) (dataChan chan string, stopChan chan bool) {
		return String2Stream(customHandler(messages))
	}, server, path)
}

func RunStreamAgentOnPath(customHandler StreamHandler, port int, path string) {
	r := gin.Default()
	RegisterStreamAgentOnPath(customHandler, r, path)
	log.Printf("Starting server on port %d", port)
	log.Fatal(r.Run(":" + strconv.Itoa(port)))
}

func RunAgentOnPath(customHandler Handler, port int, path string) {
	r := gin.Default()
	RegisterAgentOnPath(customHandler, r, path)
	log.Printf("Starting server on port %d", port)
	log.Fatal(r.Run(":" + strconv.Itoa(port)))
}

func RunAgent(customHandler Handler, port int) {
	RunAgentOnPath(customHandler, port, "/v1/chat/completions")
}

func RunStreamAgent(customHandler StreamHandler, port int) {
	RunStreamAgentOnPath(customHandler, port, "/v1/chat/completions")
}
