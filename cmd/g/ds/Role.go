package ds

type Role struct {
	Name           string   `json:"name"`
	ApplicationId  string   `json:"applicationId"`
	ChannelIds     []string `json:"channelIds"`
	ClearDelimiter string   `json:"clearDelimiter"`
	Characters     []struct {
		Desc string `json:"desc"`
	} `json:"characters"`
}
