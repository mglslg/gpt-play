package hogwarts

import (
	"fmt"
	"time"
)

func InitJokeTime() {
	ticker := time.NewTicker(3600 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("hello world")
		}
	}
}
