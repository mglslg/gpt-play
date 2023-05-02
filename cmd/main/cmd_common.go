package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/g"
)

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
