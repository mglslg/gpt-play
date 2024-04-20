package ds

// UserSession 每个用户在每个频道中都持有一个session保存其状态
type UserSession struct {
	UserChannelID   string //UserSession的唯一key
	UserId          string
	ChannelID       string
	ParentChannelID string
	UserName        string
	ClearDelimiter  string
	Model           string
	Temperature     float64
	Prompt          string
	AllowChannelIds []string //频道权限,针对非VIP用户生效
	OnConversation  bool     //是否开启上下文
	OnAt            bool     //是否需要AT机器人才会回复
}
