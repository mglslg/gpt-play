package ds

type UserSession struct {
	UserId         string  `json:"userId"`
	UserName       string  `json:"userName"`
	ClearDelimiter string  `json:"delimiter"`
	Model          string  `json:"model"`
	Temperature    float64 `json:"temperature"`
	Prompt         string  `json:"prompt"`
}
