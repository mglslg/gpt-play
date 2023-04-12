package ds

type GlobalConfig struct {
	GuildID        string `yaml:"guildId"` //公会ID(聊天室ID)
	AdminID        string `yaml:"adminId"`
	DiscordBotID   string `yaml:"discordBodId"`
	Home           string `yaml:"home"`
	MaxFetchRecord int    `yaml:"maxFetchRecord"`
	MaxUserRecord  int    `yaml:"maxUserRecord"`
}
