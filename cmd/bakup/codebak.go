package bakup

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

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

			messages = append(messages, msg)

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
