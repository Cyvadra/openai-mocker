# OpenAI Mocker
> 将任意服务包装为 OpenAI API 的标准接口格式，支持 /v1/chat/completions

## 使用方法
```golang
package main

import (
	"github.com/Cyvadra/openai-mocker/mocker"
)

func SampleReply(messages []mocker.ChatRequestMessage) string {
	return "meow"
}

func main() {
	mocker.RunAgent(SampleReply, 3001)
}

```




