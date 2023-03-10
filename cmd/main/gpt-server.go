package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Web server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/chat" {
		// 从请求中获取用户咨询的内容
		query := r.FormValue("q")
		// 设置OpenAI接口的URL和API密钥
		url := "https://api.openai.com/v1/completions"
		apikey := "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"
		// 将用户咨询的内容封装成JSON字符串
		data := fmt.Sprintf(`{"prompt": "%s", "max_tokens": 1024, "temperature": 0.7}`, query)
		// 创建新的HTTP请求
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 设置Authorization请求头部
		req.Header.Set("Authorization", "Bearer "+apikey)
		req.Header.Set("Content-Type", "application/json")
		// 发送HTTP请求
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// 读取响应数据并解析为JSON格式
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var result struct {
			Choices []struct {
				Text string `json:"text"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 将响应作为字符串返回给用户
		response := result.Choices[0].Text

		// 检查响应是否正确
		if len(result.Choices) == 0 {
			http.Error(w, "no response from OpenAI API", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, response)
	} else {
		fmt.Fprintf(w, "Hello, World!")
	}
}
