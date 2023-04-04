package ds

type GlobalConfig struct {
	ApplicationID  string `yaml:"applicationId"`
	GuildID        string `yaml:"guildId"` //公会ID(聊天室ID)
	ChannelID      string `yaml:"channelId"`
	DiscordBotID   string `yaml:"discordBodId"`
	Home           string `yaml:"home"`
	ClearCmd       string `yaml:"clearCmd"`
	Prompt         string `yaml:"prompt"`
	MaxFetchRecord int    `yaml:"maxFetchRecord"`
	MaxUserRecord  int    `yaml:"maxUserRecord"`
}
