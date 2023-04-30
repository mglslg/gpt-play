package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/g"
	"github.com/mglslg/gpt-play/cmd/g/ds"
	"github.com/mglslg/gpt-play/cmd/gpt_sdk"
	"github.com/mglslg/gpt-play/cmd/notion_sdk"
	"github.com/mglslg/gpt-play/cmd/util"
	"regexp"
	"strings"
	"time"
)

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

	//创建一忘皆空命令
	_, cmdErr := session.ApplicationCommandCreate(g.Role.ApplicationId, g.Conf.GuildID, &discordgo.ApplicationCommand{
		Name:        "一忘皆空",
		Description: "清除与" + g.Role.Name + "的聊天上下文",
	})
	if cmdErr != nil {
		logger.Fatal("create 一忘皆空 error", cmdErr)
		return nil, cmdErr
	}
	//创建导入标注内容到Notion命令
	_, cmdErr = session.ApplicationCommandCreate(g.Role.ApplicationId, g.Conf.GuildID, &discordgo.ApplicationCommand{
		Name:        "import_to_notion",
		Description: "导入标注的聊天记录到Notion",
	})
	if cmdErr != nil {
		logger.Fatal("create 导入标注到Notion error", cmdErr)
		return nil, cmdErr
	}
	session.AddHandler(onSlashCmd)

	//监听消息
	session.AddHandler(onMsgCreate)

	return session, nil
}

func onSlashCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "一忘皆空" {
		doForgetAllCmd(s, i)
	}
	if i.ApplicationCommandData().Name == "import_to_notion" {
		if doPinsToNotionCmd(s, i) {
			return
		}
	}
}

func doPinsToNotionCmd(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	// 获取钉住消息列表
	pins, err := s.ChannelMessagesPinned(i.ChannelID)
	if err != nil {
		logger.Println("Error getting pinned messages,", err)
		return true
	}
	// 获取钉住消息的内容
	for _, pin := range pins {
		client := notion_sdk.GetClient()
		notionErr := notion_sdk.AddChatHistoryEntry(client, pin.Content, time.Now())
		if notionErr != nil {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "导入失败:" + notionErr.Error(),
				},
			})
			if err != nil {
				logger.Println("Error responding to slash command: ", err)
			}
			logger.Println("Error add chat history entry to notion,", notionErr)
			return true
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "导入成功",
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
	return false
}

func doForgetAllCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		logger.Println("Error responding to slash command: ", err)
	}
}

func onMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//如果是机器人发的消息则不予理睬
	if m.Author.ID == g.Conf.DiscordBotID {
		return
	}

	if isPrivateChat(s, m) {
		if util.ContainsString(m.Author.ID, g.PrivateChatAuth.UserIds) {
			//私聊
			if m.Content == "/一忘皆空" {
				replyContent := fmt.Sprintf("%s", g.Role.ClearDelimiter)
				s.ChannelMessageSend(m.ChannelID, replyContent)
			} else {
				reply(s, m)
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "[您尚未开通私聊权限,请联系管理员Solongo]")
		}
	} else {
		if util.ContainsString(m.Author.ID, g.PrivateChatAuth.SuperUserIds) && m.Mentions != nil {
			//超级用户,不限制频道
			for _, mentioned := range m.Mentions {
				if mentioned.ID == g.Conf.DiscordBotID {
					reply(s, m)
					break
				}
			}
		} else if util.ContainsString(m.ChannelID, g.Role.ChannelIds) && m.Mentions != nil {
			//特定频道聊天,不限制用户
			for _, mentioned := range m.Mentions {
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
	conversation := geMentionContext(allMsg, m.Author.ID)
	if isPrivateChat(s, m) {
		logger.Println("/******************私聊Start:", m.Author.Username, ",privateChat:", m.Author.ID, "******************\\")
		conversation = getPrivateContext(allMsg)
	}

	//异步获取聊天记录并提示[正在输入]
	rsChnl := make(chan string)
	go callOpenAI(conversation, m.Author.Username, rsChnl)
	for {
		select {
		case gptResp := <-rsChnl:
			// Mention the user who asked the question
			msgContent := fmt.Sprintf("%s %s", m.Author.Mention(), gptResp)

			if isPrivateChat(s, m) {
				msgContent = fmt.Sprintf("%s", gptResp)
			}

			//当消息超长时拆分成两段回复用户,并且不会宕机
			var err error
			if len(msgContent) > 2000 {
				half := len(msgContent) / 2
				firstHalf := msgContent[:half]
				secondHalf := msgContent[half:]
				_, err = s.ChannelMessageSend(m.ChannelID, firstHalf)
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s %s", m.Author.Mention(), secondHalf))
			} else {
				_, err = s.ChannelMessageSend(m.ChannelID, msgContent)
			}
			if err != nil {
				logger.Println("发送discord消息失败,当前消息长度:", len(msgContent), err)
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint("[发送discord消息失败,当前消息长度:", len(msgContent), "]"))
			}

			if isPrivateChat(s, m) {
				logger.Println("\\******************私聊End:", m.Author.Username, ",privateChat:", m.Author.ID, "******************/")
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

func geMentionContext(messages []*discordgo.Message, currUserID string) *ds.Stack {
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

func getPrivateContext(messages []*discordgo.Message) *ds.Stack {
	msgStack := ds.NewStack()
	for _, msg := range messages {
		if strings.Contains(msg.Content, g.Role.ClearDelimiter) {
			return msgStack
		}
		msgStack.Push(msg)
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

	//人设
	messages = append(messages, ds.ChatMessage{
		Role:    "system",
		Content: g.Role.Characters[0].Desc,
	})

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
	if len(messages) > 20 {
		resultChannel <- abstractStrategy(messages, currUser)
	} else {
		resultChannel <- completeStrategy(messages, currUser)
	}
}

func completeStrategy(messages []ds.ChatMessage, currUser string) (resp string) {
	logger.Println("================", currUser, "================")
	for _, m := range messages {
		logger.Println(m.Role, ":", getCleanMsg(m.Content))
	}
	logger.Println("================================")
	result, _ := gpt_sdk.Chat(messages, 1)
	return result
}

func abstractStrategy(messages []ds.ChatMessage, currUser string) (resp string) {
	//处理数组越界问题
	defer func() {
		if r := recover(); r != nil {
			logger.Println("Panic occurred:", r)
		}
	}()

	lastIdx := len(messages) - 1
	lastQuestion := messages[lastIdx]

	messages[lastIdx] = ds.ChatMessage{
		Role:    "user",
		Content: "尽量详细的概括上述聊天内容",
	}
	abstract, _ := gpt_sdk.Chat(messages, 1)
	abstractMsg := make([]ds.ChatMessage, 0)

	//人设
	abstractMsg = append(abstractMsg, ds.ChatMessage{
		Role:    "system",
		Content: g.Role.Characters[0].Desc,
	})

	//上下文的概括
	abstractMsg = append(abstractMsg, ds.ChatMessage{
		Role:    "assistant",
		Content: abstract,
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

// Deprecated: 备份一下之前在倒数第二句增加的system背景
func completeStrategy_bak(messages []ds.ChatMessage, currUser string) (resp string) {
	//处理数组越界问题
	defer func() {
		if r := recover(); r != nil {
			logger.Println("Panic occurred:", r)
		}
	}()

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
