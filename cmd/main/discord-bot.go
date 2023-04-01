package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/ds"
	"github.com/mglslg/gpt-play/cmd/mygpt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Token is the token for the discord bot and chatgpt
type Token struct {
	Discord string `yaml:"discord"`
	ChatGPT string `yaml:"chatgpt"`
}

var token Token

var applicationID = "1084372136812089414"
var guildID = "1084356913816412190" //公会ID(聊天室ID)

// var channelID = "1084356914281992222" //测试频道
var channelID = "1084356913816412195"
var discordBotId = ""

var logger *log.Logger

// var home = "/Users/suolongga/app"
var home = "/app"

func main() {

	logFile := initLogger()

	readConfig()

	session, err := initDiscordSession()
	if err != nil {
		logger.Fatal("Error init discord session:", err)
		return
	}

	err = session.Open()
	if err != nil {
		logger.Fatal("Error opening connection:", err)
		return
	}

	discordBotId = session.State.User.ID

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()

	defer logFile.Close()
}

func initDiscordSession() (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + token.Discord)
	if err != nil {
		logger.Fatal("Error creating Discord session:", err)
		return nil, err
	}

	//设置机器人权限
	//intents := discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent
	intents := discordgo.IntentsAllWithoutPrivileged
	session.Identify.Intents = intents

	//创建slash命令
	_, cmdErr := session.ApplicationCommandCreate(applicationID, guildID, &discordgo.ApplicationCommand{
		Name:        "一忘皆空",
		Description: "清除与gpt机器人的聊天上下文",
	})
	if cmdErr != nil {
		logger.Fatal("create discord command error", cmdErr)
		return nil, cmdErr
	}
	session.AddHandler(onSlashCmd)

	//监听消息
	session.AddHandler(onMsgCreate)

	return session, nil
}

func initLogger() *os.File {
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("%s/deploy/logs/%s.log", home, currentDate)

	// 创建一个日志文件
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	// 创建一个日志记录器
	logger = log.New(io.MultiWriter(os.Stderr, f), "", log.LstdFlags)

	return f
}

// readConfig reads the config file and unmarshals it into the config variable
func readConfig() {
	fmt.Println("Reading config file...")

	file, err := ioutil.ReadFile(home + "/config/config.yaml")

	if err != nil {
		logger.Fatal(err.Error())
	}

	err = yaml.Unmarshal(file, &token)

	if err != nil {
		logger.Fatal(err.Error())
	}

	fmt.Println("Config file read successfully!")
}

func onSlashCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "一忘皆空" {
		//清除聊天上下文
		clearErr := clearConversation(s)
		if clearErr != nil {
			logger.Fatal("清除上下文失败", clearErr)
		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "啊天哪……我好像失忆了……",
			},
		})

		if err != nil {
			fmt.Println("Error responding to slash command: ", err)
		}
	}
}

func onMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == discordBotId {
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		logger.Fatal("Error getting channel info:", err)
		return
	}

	if channel.ID == channelID && m.Mentions != nil {
		for _, mentioned := range m.Mentions {

			logger.Println("discordBotId:", discordBotId+",mentioned.ID:", mentioned.ID, "m.Author.Mention()", m.Author.Mention())

			if mentioned.ID == discordBotId {
				mention := s.State.User.Mention()
				cleanContent := strings.Replace(m.Content, mention, "", 1)
				cleanContent = strings.TrimSpace(cleanContent)

				allMsg, e := fetchMessagesByCount(s, channelID, 30)
				if e != nil {
					logger.Fatal("抓取聊天记录失败", err)
				}

				conversation := getUserConversation(allMsg, m.Author.ID)

				aiResp, aiErr := callOpenAI(cleanContent, conversation, m.Author.Username)
				if aiErr != nil {
					logger.Fatal("Error getting response from OpenAI:", err)
					return
				}

				// Mention the user who asked the question
				userMention := m.Author.Mention()
				msgContent := fmt.Sprintf("%s %s", userMention, aiResp)

				_, err := s.ChannelMessageSend(m.ChannelID, msgContent)

				if err != nil {
					logger.Fatal("Error sending message:", err)
				}
				break
			}
		}
	}
}

func clearConversation(s *discordgo.Session) error {

	return nil
}

func getUserConversation(messages []*discordgo.Message, currUserID string) *ds.Stack {
	msgStack := ds.NewStack()
	for _, msg := range messages {
		for _, mention := range msg.Mentions {
			//找出当前用户艾特GPT以及GPT艾特当前用户的聊天记录
			if (msg.Author.ID == discordBotId && mention.ID == currUserID) || (msg.Author.ID == currUserID && mention.ID == discordBotId) {
				msgStack.Push(msg)
			}
		}
	}
	return msgStack
}

func fetchMessagesByCount(s *discordgo.Session, channelID string, count int) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message

	msgs, err := s.ChannelMessages(channelID, 100, "", "", "")

	if err != nil {
		logger.Fatal("Error fetching channel messages:", err)
		return messages, err
	}
	for index, msg := range msgs {
		if index < count {
			messages = append(messages, msg)

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
		}
	}
	return messages, nil
}

func callOpenAI(msg string, msgStack *ds.Stack, currUser string) (string, error) {

	messages := make([]ds.ChatMessage, 0)

	for !msgStack.IsEmpty() {
		msg, _ := msgStack.Pop()

		role := "user"
		if msg.Author.ID == discordBotId {
			role = "system"
		}

		messages = append(messages, ds.ChatMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	messages = append(messages, ds.ChatMessage{
		Role:    "user",
		Content: msg,
	})

	logger.Println("================", currUser, "================")
	for _, m := range messages {
		logger.Println(m)
	}
	logger.Println("================================")

	api := "https://api.openai.com/v1/chat/completions"
	payload := map[string]interface{}{
		"model":    "gpt-3.5-turbo",
		"messages": messages,
	}
	body, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(body))
	if err != nil {
		logger.Fatal("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.ChatGPT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal("Error sending request", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal("Error reading response", err)
		return "", err
	}

	chatGPTData := mygpt.ChatGPTResponse{}
	err = json.Unmarshal(body, &chatGPTData)
	if err != nil {
		logger.Fatal("Error unmarshalling response", err)
		return "", err
	}

	if len(chatGPTData.Choices) == 0 {
		return "未获取到gpt响应数据", nil
	}

	return chatGPTData.Choices[0].Message.Content, nil
}
