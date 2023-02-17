package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func StartServer() {
	url := "https://api.openai.com/v1/completions"
	token := "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"

	// 构造 HTTP 请求体
	requestBody := map[string]interface{}{
		"model":             "text-davinci-003",
		"prompt":            "我想自己在云主机上搭一个服务，有什么办法能在我提交代码之后自动在云主机上构建项目的办法吗？",
		"temperature":       0.7,
		"max_tokens":        256,
		"top_p":             1,
		"frequency_penalty": 0,
		"presence_penalty":  0,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}

	// 创建 HTTP 请求对象
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		panic(err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 发送 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 处理 HTTP 响应
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		panic(err)
	}
}
