package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/g"
	"github.com/mglslg/gpt-play/cmd/g/ds"
	"os"
)

func onBoggartSlashCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//为当前用户创建session(用户执行命令时可能早于消息监听事件,因此此处也要判断session并创建)
	us := g.GetUserSession(i.Interaction.Member.User.ID, i.Interaction.ChannelID, i.Interaction.Member.User.Username)

	if i.ApplicationCommandData().Name == "滑稽滑稽" {
		doForgetAllCmd(s, i)
	}
	if i.ApplicationCommandData().Name == "python专家" {
		onPythonExpertCmd(s, i, us)
	}
	if i.ApplicationCommandData().Name == "golang专家" {
		onGolangExpertCmd(s, i, us)
	}
	if i.ApplicationCommandData().Name == "java专家" {
		onJavaExpertCmd(s, i, us)
	}
	if i.ApplicationCommandData().Name == "linux专家" {
		onLinuxExpertCmd(s, i, us)
	}
	if i.ApplicationCommandData().Name == "网络专家" {

	}
	if i.ApplicationCommandData().Name == "自定义prompt" {

	}
	if i.ApplicationCommandData().Name == "import_to_notion" {
		if doPinsToNotionCmd(s, i) {
			return
		}
	}
}

func onPythonExpertCmd(s *discordgo.Session, i *discordgo.InteractionCreate, us *ds.UserSession) {
	us.Prompt = readPromptFromFile("python_expert")
	us.ClearDelimiter = "(博格特已变成python专家)"

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: us.ClearDelimiter,
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}

func onGolangExpertCmd(s *discordgo.Session, i *discordgo.InteractionCreate, us *ds.UserSession) {
	us.Prompt = readPromptFromFile("golang_expert")
	us.ClearDelimiter = "(博格特已变成golang专家)"

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: us.ClearDelimiter,
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}

func onJavaExpertCmd(s *discordgo.Session, i *discordgo.InteractionCreate, us *ds.UserSession) {
	us.Prompt = readPromptFromFile("java_expert")
	us.ClearDelimiter = "(博格特已变成java专家)"

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: us.ClearDelimiter,
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}

func onLinuxExpertCmd(s *discordgo.Session, i *discordgo.InteractionCreate, us *ds.UserSession) {
	us.Prompt = readPromptFromFile("linux_expert")
	us.ClearDelimiter = "(博格特已变成linux专家)"

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: us.ClearDelimiter,
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}

func readPromptFromFile(fileName string) string {
	g.Logger.Println("Reading English translator prompt file...")
	file, err := os.ReadFile("role/boggart_prompt/" + fileName)
	if err != nil {
		g.Logger.Println(err.Error())
	}
	return string(file)
}
