package notion_sdk

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	defer func(t time.Time) {
		elapsed := time.Since(t)
		fmt.Printf("myFunc took %s", elapsed)
	}(time.Now())

	client := GetClient()
	notionErr := AddChatHistoryEntry(client, "啦啦啦啦啦", time.Now())
	if notionErr != nil {
		t.Error(notionErr)
	}
}
