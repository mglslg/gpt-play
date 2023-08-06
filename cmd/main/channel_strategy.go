package main

import (
	"github.com/mglslg/gpt-play/cmd/g"
	"github.com/mglslg/gpt-play/cmd/g/ds"
)

func exeChannelStrategy(us *ds.UserSession) {
	channelId := us.ChannelID
	logger.Println("Harry>>>{}", channelId)
	if channelId == "1127815740725153812" {
		us.Prompt = g.Role.Characters[1].Desc
		us.Model = "gpt-4"
	}
}
