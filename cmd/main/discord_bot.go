package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/g"
	"github.com/mglslg/gpt-play/cmd/g/ds"
	"github.com/mglslg/gpt-play/cmd/gpt_sdk"
	"github.com/mglslg/gpt-play/cmd/util"
	"os"
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

	if g.Role.Name == "Boggart" {
		createCmd(session, "滑稽滑稽", "清除博格特上下文")
		createCmd(session, "python专家", "解答各类python相关问题")
		createCmd(session, "golang专家", "解答各类golang相关问题")
		createCmd(session, "java专家", "解答各类java相关问题")
		createCmd(session, "linux专家", "解答各类linux相关问题")
		//createCmd(session, "英文翻译", "将其它语言翻译成英文")
		//createCmd(session, "中文翻译", "将其它语言翻译成中文")
		session.AddHandler(onBoggartSlashCmd)
	} else if g.Role.Name == "Maainong" {
		//todo 暂时没有可执行的命令
	} else {
		createCmd(session, "一忘皆空", "清除与"+g.Role.Name+"的聊天上下文")
		session.AddHandler(doForgetAllCmd)
	}

	//监听消息
	session.AddHandler(onMsgCreate)

	//删除命令方法
	//session.ApplicationCommandDelete(g.Role.ApplicationId, g.Conf.GuildID, "1103997867103899679")

	return session, nil
}

func createCmd(session *discordgo.Session, cmdName string, cmdDesc string) {
	_, cmdErr := session.ApplicationCommandCreate(g.Role.ApplicationId, g.Conf.GuildID, &discordgo.ApplicationCommand{
		Name:        cmdName,
		Description: cmdDesc,
	})
	if cmdErr != nil {
		logger.Fatal("create "+cmdName+" cmd error", cmdErr)
	}
}

func onMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//如果是机器人发的消息则不予理睬
	if m.Author.ID == g.Conf.DiscordBotID {
		return
	}

	//为当前用户创建session(机器人本身也会有用户session)
	us := g.GetUserSession(m.Author.ID, m.Author.Username)

	//为userSession设置当前channelId
	us.ChannelID = m.ChannelID

	if isPrivateChat(s, us) {
		s.ChannelMessageSend(us.ChannelID, "[私聊功能暂停使用,请到聊天室里使用机器人]")
		/*if util.ContainsString(us.UserId, g.PrivateChatAuth.UserIds) {
			//私聊
			if m.Content == "/一忘皆空" {
				replyContent := fmt.Sprintf("%s", g.Role.ClearDelimiter)
				s.ChannelMessageSend(us.ChannelID, replyContent)
			} else {
				reply(s, m)
			}
		} else {
			s.ChannelMessageSend(us.ChannelID, "[您尚未开通私聊权限,请联系管理员Solongo]")
		}*/
	} else {
		if util.ContainsString(us.UserId, g.PrivateChatAuth.SuperUserIds) && m.Mentions != nil {
			//超级用户,不限制频道
			for _, mentioned := range m.Mentions {
				if mentioned.ID == g.Conf.DiscordBotID {
					reply(s, m, us)
					break
				}
			}
		} else if util.ContainsString(us.ChannelID, us.AllowChannelIds) && m.Mentions != nil {
			//特定频道聊天,不限制用户
			for _, mentioned := range m.Mentions {
				if mentioned.ID == g.Conf.DiscordBotID {
					reply(s, m, us)
					break
				}
			}
		}
	}
}

// 检查消息是否为私聊消息
func isPrivateChat(s *discordgo.Session, us *ds.UserSession) bool {
	channel, err := s.Channel(us.ChannelID)
	if err != nil {
		logger.Fatal("Error getting channel,", err)
	}
	if channel.Type == discordgo.ChannelTypeDM {
		return true
	}
	return false
}

// 回复用户消息
func reply(s *discordgo.Session, m *discordgo.MessageCreate, us *ds.UserSession) {
	allMsg, e := fetchMessagesByCount(s, us.ChannelID, g.Conf.MaxUserRecord)
	if e != nil {
		logger.Fatal("抓取聊天记录失败", e)
	}

	//获取聊天上下文
	conversation := geMentionContext(allMsg, us.UserId)

	if isPrivateChat(s, us) {
		logger.Println("/******************私聊Start:", m.Author.Username, ",privateChat:", us.UserId, "******************\\")
		conversation = getPrivateContext(allMsg)
	}

	//异步获取聊天记录并提示[正在输入]
	rsChnl := make(chan string)

	//异步调用openAI接口
	go callOpenAI(conversation, us, rsChnl)

	//异步接收接口响应
	for {
		select {
		case gptResp := <-rsChnl:
			// Mention the user who asked the question
			msgContent := fmt.Sprintf("%s %s", m.Author.Mention(), gptResp)

			if isPrivateChat(s, us) {
				msgContent = fmt.Sprintf("%s", gptResp)
			}

			//当消息超长时拆分成两段回复用户,并且不会宕机
			var err error
			if len(msgContent) > 2000 {
				half := len(msgContent) / 2
				firstHalf := msgContent[:half]
				secondHalf := msgContent[half:]
				_, err = s.ChannelMessageSend(us.ChannelID, firstHalf)
				_, err = s.ChannelMessageSend(us.ChannelID, fmt.Sprintf("%s %s", m.Author.Mention(), secondHalf))
			} else {
				_, err = s.ChannelMessageSend(us.ChannelID, msgContent)
			}
			if err != nil {
				logger.Println("发送discord消息失败,当前消息长度:", len(msgContent), err)
				_, err = s.ChannelMessageSend(us.ChannelID, fmt.Sprint("[发送discord消息失败,当前消息长度:", len(msgContent), "]"))
			}

			if isPrivateChat(s, us) {
				logger.Println("\\******************私聊End:", us.UserName, ",privateChat:", us.UserId, "******************/")
			}

			return
		default:
			s.ChannelTyping(us.ChannelID)
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

func callOpenAI(msgStack *ds.Stack, us *ds.UserSession, resultChannel chan string) {
	if msgStack.IsEmpty() {
		resultChannel <- "[没有获取到任何聊天记录,无法对话]"
		return
	}

	//翻译机器人
	if g.Role.Name == "Maainong" {
		g.Logger.Println("Reading English translator prompt file...")
		file, err := os.ReadFile("role/maainong_prompt/cn_en_translator")
		if err != nil {
			g.Logger.Println(err.Error())
		}
		translatorPrompt := string(file)
		lastMsg, _ := msgStack.GetBottomElement()
		resultChannel <- completionStrategy(getCleanMsg(lastMsg.Content), translatorPrompt, us.UserName)
		return
	}

	//打包消息列表
	messages := make([]ds.ChatMessage, 0)

	//人设
	makeSystemRole(&messages, us.Prompt)

	for !msgStack.IsEmpty() {
		msg, _ := msgStack.Pop()

		role := "user"
		if us.UserId == g.Conf.DiscordBotID {
			role = "assistant"
		}

		messages = append(messages, ds.ChatMessage{
			Role:    role,
			Content: getCleanMsg(msg.Content),
		})
	}

	//消息数大于10时使用概括策略,否则使用完整策略
	if len(messages) > 20 {
		resultChannel <- abstractChatStrategy(messages, us)
	} else {
		resultChannel <- fullChatStrategy(messages, us)
	}
}

func completionStrategy(userMessage string, prompt string, currUser string) (resp string) {
	logger.Println("================CompleteStrategy:", currUser, "================")
	logger.Println("prompt:", prompt)
	logger.Println("userMessage:", userMessage)

	result, _ := gpt_sdk.Complete(prompt, userMessage)

	logger.Println("================================")
	return result
}

func fullChatStrategy(messages []ds.ChatMessage, us *ds.UserSession) (resp string) {
	logger.Println("================", us.UserName, "================")
	for _, m := range messages {
		logger.Println(m.Role, ":", getCleanMsg(m.Content))
	}
	logger.Println("================================")

	//todo 根据channelId决定使用哪套模型
	result, _ := gpt_sdk.Chat3(messages, 0.7)

	return result
}

func abstractChatStrategy(messages []ds.ChatMessage, us *ds.UserSession) (resp string) {
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

	//todo 根据channelId决定使用哪套模型
	abstract, _ := gpt_sdk.Chat3(messages, 0)
	abstractMsg := make([]ds.ChatMessage, 0)

	//人设
	makeSystemRole(&abstractMsg, g.Role.Characters[0].Desc)

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

	logger.Println("================", us.UserName, "================")
	for _, m := range abstractMsg {
		logger.Println(m.Role, ":", getCleanMsg(m.Content))
	}
	logger.Println("================================")

	result, _ := gpt_sdk.Chat3(abstractMsg, 0.7)
	return result
}

func makeSystemRole(msg *[]ds.ChatMessage, prompt string) {
	*msg = append(*msg, ds.ChatMessage{
		Role:    "system",
		Content: prompt,
	})
}
