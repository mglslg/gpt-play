package ds

type UserSession struct {
	UserId    string `json:userId`
	UserName  string `json:userName`
	Prompt    string `json:prompt`
	Delimiter string `json:delimiter`
}
