package gpt_sdk

import (
	"bytes"
	"encoding/json"
	"github.com/mglslg/gpt-play/cmd/g"
	"github.com/mglslg/gpt-play/cmd/g/ds"
	ds2 "github.com/mglslg/gpt-play/cmd/gpt_sdk/ds"
	"io/ioutil"
	"net/http"
)

func Chat(msg []ds.ChatMessage, temperature int) (string, error) {
	api := "https://api.openai.com/v1/chat/completions"
	payload := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    msg,
		"temperature": temperature,
	}

	body, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(body))
	if err != nil {
		g.Logger.Fatal("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.SecToken.ChatGPT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.Logger.Fatal("Error sending request", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		g.Logger.Fatal("Error reading response", err)
		return "", err
	}

	chatGptResponse := ds2.ChatGPTResponse{}
	err = json.Unmarshal(body, &chatGptResponse)
	if err != nil {
		g.Logger.Fatal("Error unmarshalling response", err)
		return "", err
	}

	if len(chatGptResponse.Choices) == 0 {
		return "[未获取到gpt响应数据]", nil
	}

	return chatGptResponse.Choices[0].Message.Content, nil
}

func Complete(prompt string, temperature int) (string, error) {
	messages := make([]ds.ChatMessage, 0)

	//提示
	messages = append(messages, ds.ChatMessage{
		Role:    "system",
		Content: prompt,
	})

	api := "https://api.openai.com/v1/chat/completions"
	payload := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    messages,
		"temperature": temperature,
	}

	body, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(body))
	if err != nil {
		g.Logger.Fatal("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.SecToken.ChatGPT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.Logger.Fatal("Error sending request", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		g.Logger.Fatal("Error reading response", err)
		return "", err
	}

	chatGptResponse := ds2.ChatGPTResponse{}
	err = json.Unmarshal(body, &chatGptResponse)
	if err != nil {
		g.Logger.Fatal("Error unmarshalling response", err)
		return "", err
	}

	if len(chatGptResponse.Choices) == 0 {
		return "[未获取到gpt响应数据]", nil
	}

	return chatGptResponse.Choices[0].Message.Content, nil
}
