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
	"time"

	"github.com/bwmarrin/discordgo"
)

// Token is the token for the discord bot and chatgpt
type Token struct {
	Discord string `yaml:"discord"`
	ChatGPT string `yaml:"chatgpt"`
}

var token Token

// var channelID = "1084356914281992222" //测试频道
var channelID = "1084356913816412195"
var channelName = "gpt-play"

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

	//intents := discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent
	intents := discordgo.IntentsAllWithoutPrivileged
	dg.Identify.Intents = intents

	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	dg.AddHandler(messageCreate)

	//printChannelMsg(dg)

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

func printChannelMsg(dg *discordgo.Session) {
	startDate := time.Now().AddDate(0, 0, -1) // 1 day ago
	endDate := time.Now()

	messages, err := fetchMessagesInTimeRange(dg, channelID, startDate, endDate)
	if err != nil {
		fmt.Println("Error fetching messages: ", err)
		return
	}

	for _, message := range messages {
		fmt.Printf("%s: %s\n", message.Author.Username, message.Content)
	}
}

func fetchMessagesInTimeRange(s *discordgo.Session, channelID string, startDate, endDate time.Time) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message

	msgs, err := s.ChannelMessages(channelID, 100, "", "", "")

	if err != nil {
		fmt.Println("Error fetching channel messages: ", err)
		return messages, err
	}

	for _, msg := range msgs {
		msgTime := msg.Timestamp
		if msgTime.After(startDate) && msgTime.Before(endDate) {

			//messages = append(messages, msg)

			if msg.Content != "" {
				fmt.Printf("%s: %s\n", msg.Author.Username, msg.Content)
			}

			// 打印附件
			for _, attachment := range msg.Attachments {
				fmt.Printf("  [Attachment] %s: %s\n", attachment.Filename, attachment.URL)
			}

			// 打印嵌入内容
			for _, embed := range msg.Embeds {
				fmt.Printf("  [Embed] Title: %s, Description: %s, URL: %s\n", embed.Title, embed.Description, embed.URL)
			}

			// 打印自定义表情
			for _, reaction := range msg.Reactions {
				fmt.Printf("  [Reaction] Emoji: %s, Count: %d\n", reaction.Emoji.Name, reaction.Count)
			}
		} else if msgTime.Before(startDate) {
			return messages, nil
		}
	}

	return messages, nil
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

	if channel.Name == channelName && m.Mentions != nil {
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

				// Mention the user who asked the question
				userMention := m.Author.Mention()
				msgContent := fmt.Sprintf("%s %s", userMention, aiResp)

				_, err := s.ChannelMessageSend(m.ChannelID, msgContent)

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
