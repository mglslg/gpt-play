package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/notion_sdk"
	"time"
)

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
