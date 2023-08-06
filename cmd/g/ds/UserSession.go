package ds

// UserSession 每个用户在每个频道中都持有一个session保存其状态
type UserSession struct {
	UserChannelID   string //UserSession的唯一key
	UserId          string
	ChannelID       string
	UserName        string
	ClearDelimiter  string
	Model           string
	Temperature     float64
	Prompt          string
	AllowChannelIds []string
}
