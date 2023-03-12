// Package main
package main

//皓哥写的代码,瑞思拜
import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Token is the token for the discord bot and chatgpt
type Token struct {
	Discord string `yaml:"discord"`
	ChatGPT string `yaml:"chatgpt"`
}

var token Token

// ReadConfig reads the config file and unmarshals it into the config variable
func ReadConfig() error {
	fmt.Println("Reading config file...")

	configFilePath := "config/config.yaml"
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file, err := ioutil.ReadFile(filepath.Join(workingDir, configFilePath))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = yaml.Unmarshal(file, &token)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Config file read successfully!")

	return nil

}

func Start() {
	// 创建 ServeMux 实例
	mux := http.NewServeMux()

	// 注册路由及其对应的处理程序
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello, World!")
		fmt.Fprintf(w, "Hello, World!")
	})

	mux.HandleFunc("/chat", messageHandler)

	fmt.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	ChatGPTResponse, err := callChatGPT("默写锄禾日当午")
	if err != nil {
		log.Println(err.Error())
		fmt.Println(err.Error())
		return
	}
	fmt.Println(ChatGPTResponse)
	log.Print(ChatGPTResponse)
}

// ChatGPTResponse is the response from the chatgpt api
type ChatGPTResponse struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Model  string `json:"model"`
	Usage  struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func callChatGPT(msg string) (string, error) {
	api := "https://api.openai.com/v1/chat/completions"
	body := []byte(`{
  		"model": "text-davinci-003",
  		"prompt": "` + msg + `",
  		"temperature": 0.7,
  		"max_tokens": 256,
		"top_p": 1,
  		"frequency_penalty": 0,
  		"presence_penalty": 0
	}`)

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.ChatGPT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response", err)
		return "", err
	}

	log.Println(body)

	chatGPTData := ChatGPTResponse{}
	err = json.Unmarshal(body, &chatGPTData)
	if err != nil {
		fmt.Println("Error unmarshalling response", err)
		return "", err
	}
	return chatGPTData.Choices[0].Message.Content, nil
}

// JSONEscape escape the string
func JSONEscape(str string) string {
	b, err := json.Marshal(str)
	if err != nil {
		return str
	}
	s := string(b)
	return s[1 : len(s)-1]
}

func main() {
	err := ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Start()
}
