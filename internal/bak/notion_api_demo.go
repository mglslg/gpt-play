package bak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	NotionAPIKey  = "YOUR_NOTION_API_KEY"
	DatabaseID    = "69e551fe013249f6a0c6d0d7d7cbba9f"
	NotionAPIURL  = "https://api.notion.com/v1/pages"
	NotionVersion = "2021-08-16"
)

type Text struct {
	Content string `json:"content"`
}

type Title struct {
	Title []Text `json:"title"`
}

type RichText struct {
	RichText []Text `json:"rich_text"`
}

type Number struct {
	Number int `json:"number"`
}

type Date struct {
	Start string `json:"start"`
}

type Select struct {
	Name string `json:"name"`
}

type NotionProperty struct {
	TTTitle *Title    `json:"TTTitle,omitempty"`
	Desc    *RichText `json:"Desc,omitempty"`
	Age     *Number   `json:"Age,omitempty"`
	DDDate  *Date     `json:"DDDate,omitempty"`
	City    *Select   `json:"City,omitempty"`
}

type Parent struct {
	DatabaseID string `json:"database_id"`
}

type NotionPage struct {
	Parent     Parent         `json:"parent"`
	Properties NotionProperty `json:"properties"`
}

func main() {
	page := NotionPage{
		Parent: Parent{
			DatabaseID: DatabaseID,
		},
		Properties: NotionProperty{
			TTTitle: &Title{
				Title: []Text{
					{
						Content: "测试标题",
					},
				},
			},
			Desc: &RichText{
				RichText: []Text{
					{
						Content: "这是描述",
					},
				},
			},
			Age: &Number{
				Number: 28,
			},
			DDDate: &Date{
				Start: "2023-04-13",
			},
			City: &Select{
				Name: "上海",
			},
		},
	}

	jsonData, err := json.Marshal(page)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", NotionAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", NotionVersion)
	req.Header.Set("Authorization", "Bearer "+NotionAPIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		return
	}

	fmt.Println("Page created successfully")
}
