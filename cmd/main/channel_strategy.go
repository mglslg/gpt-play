package main

import (
	"github.com/mglslg/gpt-play/cmd/g"
	"github.com/mglslg/gpt-play/cmd/g/ds"
	"github.com/mglslg/gpt-play/cmd/util"
)

func setChannelStatus(us *ds.UserSession) {
	channelId := us.ChannelID

	//gpt-4
	if channelId == "1127815740725153812" {
		us.Prompt = g.Role.Characters[1].Desc
		us.Model = "gpt-4"
	}
	//translate
	if channelId == "1095947683597914162" {
		us.OnConversation = false
		us.OnAt = false
	}
}

// 超级用户或特定频道才有权限触发机器人回复
func hasChannelPrivilege(us *ds.UserSession) bool {
	if util.ContainsString(us.UserId, g.PrivateChatAuth.SuperUserIds) {
		return true
	}
	if util.ContainsString(us.ChannelID, us.AllowChannelIds) {
		return true
	}
	return false
}
