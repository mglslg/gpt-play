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

func Chat3(msg []ds.ChatMessage, temperature float64) (string, error) {
	return Chat(msg, temperature, "gpt-3.5-turbo")
}

func Chat4(msg []ds.ChatMessage, temperature float64) (string, error) {
	return Chat(msg, temperature, "gpt-4")
}

// todo 最终还是要使用这个,model从跟外面的session传入进来
func Chat(msg []ds.ChatMessage, temperature float64, model string) (string, error) {
	api := "https://api.openai.com/v1/chat/completions"
	payload := map[string]interface{}{
		"model":       model,
		"messages":    msg,
		"temperature": temperature,
	}

	body, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(body))
	if err != nil {
		g.Logger.Println("Error creating request:", err)
		return "[Error creating request:" + err.Error() + "]", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.SecToken.ChatGPT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.Logger.Println("Error sending request", err)
		return "[Error sending request:" + err.Error() + "]", err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		g.Logger.Println("Error reading response", err)
		return "[Error reading response:" + err.Error() + "]", err
	}

	chatGptResponse := ds2.ChatGPTResponse{}
	err = json.Unmarshal(body, &chatGptResponse)
	if err != nil {
		g.Logger.Println("[Error unmarshalling response]", err)
		return "[Error unmarshalling response:" + err.Error() + "]", err
	}

	if len(chatGptResponse.Choices) == 0 {
		return "[未获取到gpt响应数据]", nil
	}
	g.Logger.Println(">>>>>gpt响应:", chatGptResponse.Choices[0].Message.Content)
	g.Logger.Println(">>>>>finish原因:", chatGptResponse.Choices[0].FinishReason)
	g.Logger.Println(">>>>>已花费token:", chatGptResponse.Usage.TotalTokens)

	return chatGptResponse.Choices[0].Message.Content, nil
}

func Complete(prompt string, message string, temperature int, model string) (string, error) {
	if model == "" {
		model = "text-davinci-003"
	}

	prompt = prompt + "```" + message + "```"

	api := "https://api.openai.com/v1/completions"
	payload := map[string]interface{}{
		"model":       model,
		"prompt":      prompt,
		"temperature": temperature,
		"max_tokens":  2048,
		"n":           1,
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

	g.Logger.Println(">>>>>prompt:", prompt)
	g.Logger.Println(">>>>>gpt响应:", chatGptResponse.Choices[0].Text)
	g.Logger.Println(">>>>>finish原因:", chatGptResponse.Choices[0].FinishReason)
	g.Logger.Println(">>>>>已花费token:", chatGptResponse.Usage.TotalTokens)

	return chatGptResponse.Choices[0].Text, nil
}
