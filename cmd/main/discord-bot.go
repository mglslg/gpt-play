package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mglslg/gpt-play/cmd/bakup"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
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
	//file, err := ioutil.ReadFile("/Users/suolongga/app/config/config.yaml")
	file, err := ioutil.ReadFile("/app/config/config.yaml")

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

func main() {
	err := ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dg, err := discordgo.New("Bot " + token.Discord)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel info:", err)
		return
	}

	if channel.Name == "gpt-play" && m.Mentions != nil {
		for _, user := range m.Mentions {
			if user.ID == s.State.User.ID {
				mention := s.State.User.Mention()
				cleanContent := strings.Replace(m.Content, mention, "", 1)
				cleanContent = strings.TrimSpace(cleanContent)

				aiResp, aiErr := callOpenAI(cleanContent)
				if aiErr != nil {
					fmt.Println("Error getting response from OpenAI:", err)
					return
				}

				response := fmt.Sprintf("%s", aiResp)
				_, err := s.ChannelMessageSend(m.ChannelID, response)
				if err != nil {
					fmt.Println("Error sending message:", err)
				}
				break
			}
		}
	}
}

func callOpenAI(msg string) (string, error) {
	api := "https://api.openai.com/v1/chat/completions"
	body := []byte(`{
		"model": "gpt-3.5-turbo",
		"messages": [
		  {
			"role": "user",
			"content": "` + bakup.JSONEscape(msg) + `"
		  }
		]
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

	chatGPTData := bakup.ChatGPTResponse{}
	err = json.Unmarshal(body, &chatGPTData)
	if err != nil {
		fmt.Println("Error unmarshalling response", err)
		return "", err
	}
	return chatGPTData.Choices[0].Message.Content, nil
}
