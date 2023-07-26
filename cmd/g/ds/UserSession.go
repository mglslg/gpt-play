package ds

type UserSession struct {
	UserId          string
	UserName        string
	ClearDelimiter  string
	Model           string
	Temperature     float64
	Prompt          string
	ChannelID       string
	AllowChannelIds []string
}
