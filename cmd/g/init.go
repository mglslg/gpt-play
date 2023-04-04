package g

import (
	"fmt"
	"github.com/mglslg/gpt-play/cmd/ds"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var Logger *log.Logger
var Conf ds.GlobalConfig
var SecToken ds.Token

// readConfig reads the config file and unmarshals it into the config variable
func InitConfig() {
	fmt.Println("Reading config file...")

	file, err := ioutil.ReadFile("config/config.yaml")

	if err != nil {
		fmt.Println("Reading config file failed!", err)
		return
	}

	err = yaml.Unmarshal(file, &Conf)

	if err != nil {
		fmt.Println("Resolve config file failed!", err)
		return
	}

	fmt.Println("Config file read successfully!")
}

func InitLogger() *os.File {
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("%s/deploy/logs/%s.log", Conf.Home, currentDate)

	// 创建一个日志文件
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	// 创建一个日志记录器
	Logger = log.New(io.MultiWriter(os.Stderr, f), "", log.LstdFlags)

	return f
}

func InitSecretConfig() {
	fmt.Println("Reading secret config file...")

	file, err := ioutil.ReadFile(Conf.Home + "/config/config.yaml")

	if err != nil {
		Logger.Fatal(err.Error())
	}

	err = yaml.Unmarshal(file, &SecToken)

	if err != nil {
		Logger.Fatal(err.Error())
	}

	Logger.Println("Secret Config file read successfully!")
}
