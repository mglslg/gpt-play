package ds

type Role struct {
	Name           string   `json:"name"`
	ApplicationId  string   `json:"applicationId"`
	ChannelIds     []string `json:"channelIds"`
	ClearDelimiter string   `json:"clearDelimiter"`
	Temperature    float64  `json:"temperature""`
	Characters     []struct {
		Desc string `json:"desc"`
	} `json:"characters"`
}
