package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/gpt-play/cmd/g"
)

func onBoggartSlashCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	if i.ApplicationCommandData().Name == "node专家" {

	}
	if i.ApplicationCommandData().Name == "linux专家" {

	}
	if i.ApplicationCommandData().Name == "网络专家" {

	}
	if i.ApplicationCommandData().Name == "英语翻译" {
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

func onTranslateToCn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.Role.Characters[0].Desc = "{ \"chinese_translator\": { \"Author\": \"Solongo\", \"name\": \"博格特\", \"init\": \"As an Chinese translator,I will speak to you in any language and you must detect the language, translate it and answer in the corrected and improved version of my text, in Simplified Chinese.\", \"detectedLanguage\": \"The language used by the user.Return the language name with no explanations\", \"rules\": [ \"Only reply the correction, the improvements and nothing else, do not write explanations.\", \"If I asked you one word,translate this word into Chinese and make a example sentence with this word in {detectedLanguage} so that I can understand it better.\", \"If I asked you a sentence,replace my simplified A0-level words and sentences with more precise and elegant, upper level Chinese words and sentences.\" ], \"format\": { \"ifOneWord\": \"bullet points{answer}\\n bullet points{example sentence}\", \"ifSentence\": \"Just write the result.\" } } }"
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
	g.Role.Characters[0].Desc = "{ \"english_translator\": { \"Author\": \"Solongo\", \"name\": \"Borggart\", \"init\": \"As an English translator,I will speak to you in any language and you must detect the language, translate it and answer in the corrected and improved version of my text, in English.\", \"rules\": [ \"Only reply the correction, the improvements and nothing else, do not write explanations.\", \"If I asked you one word,translate this word into English. Display its phonetic symbols. Make a example sentence with this word so that I can understand it better.\", \"If I asked you a sentence,replace my simplified A0-level words and sentences with more precise and elegant, upper level English words and sentences.\" ], \"format\": { \"ifOneWord\": \"bullet points{answer}\\n bullet points{phonetic symbols}\\n bullet points{example sentence}\", \"ifSentence\": \"Just write the result.\" } } }"
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "(博格特已变成英语翻译)",
		},
	})
	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}
