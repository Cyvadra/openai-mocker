# OpenAI Mocker
> 将任意服务包装为 OpenAI API 的标准接口格式，支持 /v1/chat/completions

## 使用方法
```golang

package main

import (
	"github.com/Cyvadra/openai-mocker/mocker"
)

func main() {
  // change handler to your own function
	mocker.RunAgent(mocker.SampleReply, 3001)
}

```




