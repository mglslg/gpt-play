package main

import (
	"fmt"
	"github.com/mglslg/gpt-play/internal"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// 将 "/gpt-play" 注册为全局服务路径
	gptPlayMux := http.NewServeMux()
	gptPlayMux.HandleFunc("/hello", helloHandler)

	mux.Handle("/gpt-play/", http.StripPrefix("/gpt-play", gptPlayMux))

	fmt.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}

	internal.StartServer()
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := fmt.Fprint(w, "hello world")
		if err != nil {
			return
		}
	} else {
		http.Error(w, "Post hello world", http.StatusMethodNotAllowed)
	}
}
