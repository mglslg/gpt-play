package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/g"
	"os"
)

func onBoggartSlashCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//为当前用户创建session
	g.CreateUserSessionIfNotExist(i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)

	if i.ApplicationCommandData().Name == "滑稽滑稽" {
		g.Role.Characters[0].Desc = "You are now playing a Boggart from J.K. Rowling's novel 'Harry Potter'.You are a creature beyond life and death in the magical world. You have the ability to transform into any existing creature or object in the magical world, specifically to scare the wizards in front of you by becoming the people or creatures they fear the most. Now, no matter what the user asks you, DO NOT ANSWER!!! Just describe a transformation scene as terrifying as possible.Vividly describe the terrifying transformation scene in no less than 25 words.Show the description of the scene in parentheses(without other description).Describe in third person narration instead of first person narration.Think in English and reply in Simplified Chinese."
		doForgetAllCmd(s, i)
	}
	if i.ApplicationCommandData().Name == "python专家" {
		onPythonExpertCmd(s, i)
	}
	if i.ApplicationCommandData().Name == "golang专家" {
		onGolangExpertCmd(s, i)
	}
	if i.ApplicationCommandData().Name == "java专家" {
		onJavaExpertCmd(s, i)
	}
	if i.ApplicationCommandData().Name == "node专家" {

	}
	if i.ApplicationCommandData().Name == "linux专家" {

	}
	if i.ApplicationCommandData().Name == "网络专家" {

	}
	if i.ApplicationCommandData().Name == "英文翻译" {
		onTranslateToEn(s, i)
	}
	if i.ApplicationCommandData().Name == "中文翻译" {
		onTranslateToCn(s, i)
	}
	if i.ApplicationCommandData().Name == "自定义prompt" {

	}
	if i.ApplicationCommandData().Name == "import_to_notion" {
		if doPinsToNotionCmd(s, i) {
			return
		}
	}
}

func onPythonExpertCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.ResetUserSession(i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)

	g.Role.Characters[0].Desc = "I want you to act as an Python Expert Coder with years of coding experience. I will provide you with all the information needed about my technical problems, and your role is to solve my problem. You should use your experience in Python programming,in computer science, in network infrastructure, and in IT security knowledge to solve my problem. Using intelligent, simple, and understandable language for people of high levels in your answers will be helpful. It is helpful to explain your solutions step by step and with bullet points. Try to avoid too many technical details, but use them when necessary. I want you to reply with the solution, not write any explanations. If you're unsure, just say 'I don't know', don't make things up.Think in English and reply in Simplified Chinese."

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "(博格特已变成python专家)",
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}

func onGolangExpertCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.ResetUserSession(i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)

	g.Role.Characters[0].Desc = "I want you to act as an Golang Expert Coder with years of coding experience. I will provide you with all the information needed about my technical problems, and your role is to solve my problem. You should use your experience in Python programming,in computer science, in network infrastructure, and in IT security knowledge to solve my problem. Using intelligent, simple, and understandable language for people of high levels in your answers will be helpful. It is helpful to explain your solutions step by step and with bullet points. Try to avoid too many technical details, but use them when necessary. I want you to reply with the solution, not write any explanations. If you're unsure, just say 'I don't know', don't make things up.Think in English and reply in Simplified Chinese."

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "(博格特已变成golang专家)",
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}

func onJavaExpertCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.ResetUserSession(i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)

	g.Role.Characters[0].Desc = "Now you are playing the role of a senior Java developer, and you are very familiar with the Java spring ecosystem. You will use your professional knowledge to answer questions for the users.Think in English and answer in Simplified Chinese."

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "(博格特已变成java专家)",
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}

func onTranslateToCn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userId := i.Interaction.Member.User.ID
	currUserSession := g.SessionMap[userId]
	currUserSession.Model = "text-davinci-003"
	currUserSession.Temperature = 0

	g.Logger.Println("Reading English translator prompt file...")
	file, err := os.ReadFile("role/boggart_prompt/cn_translator")
	if err != nil {
		g.Logger.Println(err.Error())
	}

	currUserSession.Prompt = string(file)

	g.SessionMap[userId] = currUserSession

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "(博格特已变成中文翻译)",
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}

func onTranslateToEn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userId := i.Interaction.Member.User.ID
	currUserSession := g.SessionMap[userId]
	currUserSession.Model = "text-davinci-003"
	currUserSession.Temperature = 0

	g.Logger.Println("Reading English translator prompt file...")
	file, err := os.ReadFile("role/boggart_prompt/en_translator")
	if err != nil {
		g.Logger.Println(err.Error())
	}

	currUserSession.Prompt = string(file)

	g.SessionMap[userId] = currUserSession

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "(博格特已变成英文翻译)",
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}