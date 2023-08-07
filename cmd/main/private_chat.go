package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/g"
	"github.com/mglslg/gpt-play/cmd/g/ds"
	"strings"
)

// Deprecated 检查消息是否为私聊消息
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

// Deprecated
func privateReply(s *discordgo.Session, m *discordgo.MessageCreate, us *ds.UserSession) {
	//allMsg, e := fetchMessagesByCount(s, us.ChannelID, g.Conf.MaxUserRecord)
	//if e != nil {
	//	logger.Fatal("抓取聊天记录失败", e)
	//}
	//
	////获取聊天上下文
	//conversation := geMentionContext(allMsg, us)
	//
	//if isPrivateChat(s, us) {
	//	logger.Println("/******************私聊Start:", m.Author.Username, ",privateChat:", us.UserId, "******************\\")
	//	conversation = getPrivateContext(allMsg)
	//}
}

// Deprecated
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
