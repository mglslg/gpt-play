package ds

// UserSession Boggart对应的用户会话,userId为当前用户的ID,其余字段为当前Boggart所使用的模型的相关信息
type UserSession struct {
	UserId         string  `json:"userId"`
	UserName       string  `json:"userName"`
	ClearDelimiter string  `json:"delimiter"`
	Model          string  `json:"model"`
	Temperature    float64 `json:"temperature"`
	Prompt         string  `json:"prompt"`
	ChannelID      string  `json:"channelId"`
}
