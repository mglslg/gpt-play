package main

import (
	"flag"
	"fmt"
	"github.com/mglslg/gpt-play/cmd/g"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger *log.Logger

func main() {
	var roleName string
	var configFilePath string

	//本地调试按需修改
	flag.StringVar(&roleName, "role", "Maainong", "The role of the bot")
	flag.StringVar(&configFilePath, "config", "config/home_config.yaml", "path to config file")
	//flag.StringVar(&configFilePath, "config", "config/company_config.yaml", "path to config file")
	flag.Parse()

	g.InitConfig(configFilePath)
	g.InitRole(roleName)
	logFile := g.InitLogger()
	logger = g.Logger
	g.InitSecretConfig()
	g.InitPrivateChatAuth()
	g.InitSessionMap()

	session, err := initDiscordSession()

	if err != nil {
		logger.Fatal("Error g discord session:", err)
		return
	} else {
		logger.Println("Session init successfully")
	}

	err = session.Open()
	if err != nil {
		logger.Fatal("Error opening connection:", err)
		return
	}

	g.Conf.DiscordBotID = session.State.User.ID

	logger.Println("Bot is now running.")
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()

	defer logFile.Close()
}
