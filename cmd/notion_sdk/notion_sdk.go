package notion_sdk

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	notionAPIBaseURL = "https://api.notion.com/v1"
	notionAPIVersion = "2022-06-07"
	databaseID       = "69e551fe013249f6a0c6d0d7d7cbba9f"
	apiKey           = "secret_QHiT0ONKwhJwSbVA61032yWr0M5PK2OSZvicnkGk3VJ"
)

type Page struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type CreatePageResponse struct {
	Page
}

func importToNotion() {
	client := resty.New()

	// Add a new entry to the chat_history database
	err := addChatHistoryEntry(client, "Hello, Notion!", time.Now())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding entry to chat_history database: %v\n", err)
		os.Exit(1)
	}
}

func addChatHistoryEntry(client *resty.Client, title string, date time.Time) error {
	// Prepare the request payload
	payload := map[string]interface{}{
		"parent": map[string]string{
			"database_id": databaseID,
		},
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"title": []interface{}{},
			},
			"date": map[string]interface{}{
				"date": map[string]interface{}{
					"start": date.Format(time.RFC3339),
				},
			},
		},
	}

	// Create the title text block
	payload["properties"].(map[string]interface{})["title"].(map[string]interface{})["title"] = []interface{}{
		map[string]interface{}{
			"type": "text",
			"text": map[string]interface{}{
				"content": title,
			},
		},
	}

	// Send the request to the Notion API
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Notion-Version", notionAPIVersion).
		SetHeader("Authorization", "Bearer "+apiKey).
		SetBody(payload).
		Post(notionAPIBaseURL + "/pages")

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("request returned an error: %v", resp.Status())
	}

	var createPageResponse CreatePageResponse
	if err := json.Unmarshal(resp.Body(), &createPageResponse); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Printf("Added chat_history entry with ID %s\n", createPageResponse.ID)

	return nil
}
