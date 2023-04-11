package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/g"
	"github.com/mglslg/gpt-play/cmd/g/ds"
	"github.com/mglslg/gpt-play/cmd/gpt_sdk"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var logger *log.Logger
var adminId = "1049923686288863283"
var irmuunId = "1094082490500657312"

func main() {
	//默认使用Snape
	var roleName string
	flag.StringVar(&roleName, "role", "Hermione", "The role of the bot")

	//默认使用local_config.yaml
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config/local_config.yaml", "path to config file")

	flag.Parse()

	g.InitConfig(configFilePath)
	g.InitRole(roleName)
	logFile := g.InitLogger()
	logger = g.Logger
	g.InitSecretConfig()

	session, err := initDiscordSession()

	if err != nil {
		logger.Fatal("Error g discord session:", err)
		return
	}

	err = session.Open()
	if err != nil {
		logger.Fatal("Error opening connection:", err)
		return
	}

	g.Conf.DiscordBotID = session.State.User.ID

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()

	defer logFile.Close()
}

func initDiscordSession() (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + g.SecToken.Discord)
	if err != nil {
		logger.Fatal("Error creating Discord session:", err)
		return nil, err
	}

	//设置机器人权限
	//intents := discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent
	intents := discordgo.IntentsAllWithoutPrivileged
	session.Identify.Intents = intents

	//创建slash命令
	_, cmdErr := session.ApplicationCommandCreate(g.Role.ApplicationId, g.Conf.GuildID, &discordgo.ApplicationCommand{
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

func onSlashCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "一忘皆空" {
		//清除聊天上下文(实际上就是打印一句话,后面取聊天记录时按照它作分隔)

		userMention := i.Member.User.Mention()
		replyContent := fmt.Sprintf("%s %s", userMention, g.Role.ClearDelimiter)

		//这边需要改成异步的，是否改成异步就可以呢？
		//aiCreate, cerro := gpt_sdk.Complete("你现在扮演JK罗琳的小说《哈利波特》里的斯内普教授，你有一头油腻的黑发，穿着一身黑袍，是一个非常优秀的巫师，会很多厉害的咒语。然而你的人品和脾气都非常差，为人傲慢刻薄残忍，缺乏耐心和同情心，折磨别人是你的乐趣。此时此刻与你聊天的人正被你激怒，他飞快的举起魔杖对着你大喊「一忘皆空」。你大声咒骂并拔出魔杖，但为时已晚，仍然不幸被咒语击中，失去了一段记忆。现在你发挥想象，把这段场景用生动的语言描绘出来，要突出斯内普的咒骂以及被咒语击中的瞬间，不超过100个token。", 1)
		//if cerro != nil {
		//	logger.Fatal("获取一忘皆空创作失败", cerro)
		//}
		//logger.Println(aiCreate)
		//replyContent = aiCreate

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: replyContent,
			},
		})

		if err != nil {
			fmt.Println("Error responding to slash command: ", err)
		}
	}
}

func onMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//如果是机器人发的消息则不予理睬
	if m.Author.ID == g.Conf.DiscordBotID {
		return
	}

	if isPrivateChat(s, m) && m.Author.ID == adminId {
		if m.Content == "/一忘皆空" {
			replyContent := fmt.Sprintf("%s", g.Role.ClearDelimiter)
			s.ChannelMessageSend(m.ChannelID, replyContent)
		} else {
			reply(s, m)
		}
	} else if (m.Author.ID == adminId || m.Author.ID == irmuunId) && m.Mentions != nil {
		for _, mentioned := range m.Mentions {
			logger.Println("discordBotId:", g.Conf.DiscordBotID+",mentioned.ID:", mentioned.ID)
			if mentioned.ID == g.Conf.DiscordBotID {
				reply(s, m)
				break
			}
		}
	} else {
		if m.ChannelID == g.Role.ChannelIds[0] && m.Mentions != nil {
			for _, mentioned := range m.Mentions {
				logger.Println("discordBotId:", g.Conf.DiscordBotID+",mentioned.ID:", mentioned.ID)
				if mentioned.ID == g.Conf.DiscordBotID {
					reply(s, m)
					break
				}
			}
		}
	}
}

// 检查消息是否为私聊消息
func isPrivateChat(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		logger.Fatal("Error getting channel,", err)
	}
	if channel.Type == discordgo.ChannelTypeDM {
		return true
	}
	return false
}

// 回复用户消息
func reply(s *discordgo.Session, m *discordgo.MessageCreate) {
	allMsg, e := fetchMessagesByCount(s, m.ChannelID, g.Conf.MaxUserRecord)
	if e != nil {
		logger.Fatal("抓取聊天记录失败", e)
	}

	//获取聊天上下文
	conversation := getPublicContext(allMsg, m.Author.ID)
	if isPrivateChat(s, m) {
		conversation = getPrivateContext(allMsg, m.Author.ID)
	}

	//异步获取聊天记录并提示[正在输入]
	rsChnl := make(chan string)
	go callOpenAI(conversation, m.Author.Username, rsChnl)
	for {
		select {
		case gptResp := <-rsChnl:
			// Mention the user who asked the question
			msgContent := fmt.Sprintf("%s %s", m.Author.Mention(), gptResp)

			_, err := s.ChannelMessageSend(m.ChannelID, msgContent)

			if err != nil {
				logger.Fatal("Error sending message:", err)
			}
			return
		default:
			s.ChannelTyping(m.ChannelID)
			time.Sleep(5 * time.Second)
		}
	}
}

func getCleanMsg(content string) string {
	// 创建一个正则表达式，用于匹配尖括号及其内容，格式为：<@数字>
	re := regexp.MustCompile(`<@(\d+)>`)

	// 使用正则表达式替换匹配的内容为空字符串
	cleanedMsg := re.ReplaceAllString(content, "")

	return cleanedMsg
}

func getPublicContext(messages []*discordgo.Message, currUserID string) *ds.Stack {
	msgStack := ds.NewStack()
	for _, msg := range messages {
		for _, mention := range msg.Mentions {
			//找出当前用户艾特GPT以及GPT艾特当前用户的聊天记录
			if (msg.Author.ID == g.Conf.DiscordBotID && mention.ID == currUserID) || (msg.Author.ID == currUserID && mention.ID == g.Conf.DiscordBotID) {
				//一旦发现clear命令的分隔符则直接终止向消息栈push,直接返回
				if strings.Contains(msg.Content, g.Role.ClearDelimiter) {
					return msgStack
				}
				msgStack.Push(msg)
			}
		}
	}
	return msgStack
}

func getPrivateContext(messages []*discordgo.Message, currUserID string) *ds.Stack {

	//todo 待实现
	msgStack := ds.NewStack()
	for _, msg := range messages {
		for _, mention := range msg.Mentions {
			//找出当前用户艾特GPT以及GPT艾特当前用户的聊天记录
			if (msg.Author.ID == g.Conf.DiscordBotID && mention.ID == currUserID) || (msg.Author.ID == currUserID && mention.ID == g.Conf.DiscordBotID) {
				//一旦发现clear命令的分隔符则直接终止向消息栈push,直接返回
				if strings.Contains(msg.Content, g.Role.ClearDelimiter) {
					return msgStack
				}
				msgStack.Push(msg)
			}
		}
	}
	return msgStack
}

func fetchMessagesByCount(s *discordgo.Session, channelID string, count int) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message

	msgs, err := s.ChannelMessages(channelID, g.Conf.MaxFetchRecord, "", "", "")

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

func callOpenAI(msgStack *ds.Stack, currUser string, resultChannel chan string) {
	if msgStack.IsEmpty() {
		resultChannel <- "[没有获取到任何聊天记录,无法对话]"
		return
	}

	//打包消息列表
	messages := make([]ds.ChatMessage, 0)
	for !msgStack.IsEmpty() {
		msg, _ := msgStack.Pop()

		role := "user"
		if msg.Author.ID == g.Conf.DiscordBotID {
			role = "assistant"
		}

		messages = append(messages, ds.ChatMessage{
			Role:    role,
			Content: getCleanMsg(msg.Content),
		})
	}

	//消息数大于10时使用概括策略,否则使用完整策略
	if len(messages) > 10 {
		resultChannel <- abstractStrategy(messages, currUser)
	} else {
		resultChannel <- completeStrategy(messages, currUser)
	}
}

func completeStrategy(messages []ds.ChatMessage, currUser string) (resp string) {
	lastIdx := len(messages) - 1
	lastQuestion := messages[lastIdx]

	//给倒数第二条聊天记录设置人设，降低逃逸概率
	messages[lastIdx] = ds.ChatMessage{
		Role:    "system",
		Content: g.Role.Characters[0].Desc,
	}
	messages = append(messages, lastQuestion)

	logger.Println("================", currUser, "================")
	for _, m := range messages {
		logger.Println(m.Role, ":", getCleanMsg(m.Content))
	}
	logger.Println("================================")
	result, _ := gpt_sdk.Chat(messages, 1)
	return result
}

func abstractStrategy(messages []ds.ChatMessage, currUser string) (resp string) {
	lastIdx := len(messages) - 1
	lastQuestion := messages[lastIdx]

	messages[lastIdx] = ds.ChatMessage{
		Role:    "user",
		Content: "尽量详细的概括上述聊天内容",
	}
	abstract, _ := gpt_sdk.Chat(messages, 1)
	abstractMsg := make([]ds.ChatMessage, 0)

	//上下文的概括
	abstractMsg = append(abstractMsg, ds.ChatMessage{
		Role:    "assistant",
		Content: abstract,
	})

	//人设
	abstractMsg = append(abstractMsg, ds.ChatMessage{
		Role:    "system",
		Content: g.Role.Characters[0].Desc,
	})

	//用户问题
	abstractMsg = append(abstractMsg, ds.ChatMessage{
		Role:    "user",
		Content: lastQuestion.Content,
	})

	logger.Println("================", currUser, "================")
	for _, m := range abstractMsg {
		logger.Println(m.Role, ":", getCleanMsg(m.Content))
	}
	logger.Println("================================")

	result, _ := gpt_sdk.Chat(abstractMsg, 1)
	return result
}
