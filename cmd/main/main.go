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
	//默认使用多比
	var roleName string
	flag.StringVar(&roleName, "role", "Dobby", "The role of the bot")
	//默认使用local_config.yaml
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config/local_config.yaml", "path to config file")
	flag.Parse()

	g.InitConfig(configFilePath)
	g.InitRole(roleName)
	logFile := g.InitLogger()
	logger = g.Logger
	g.InitSecretConfig()
	g.InitPrivateChatAuth()

	session, err := initDiscordSession()

	if err != nil {
		logger.Fatal("Error g discord session:", err)
		return
	}

	err = session.Open()
	if err != nil {
		logger.Fatal("Error opening connection:", err)
		return
	}

	g.Conf.DiscordBotID = session.State.User.ID

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()

	defer logFile.Close()
}
