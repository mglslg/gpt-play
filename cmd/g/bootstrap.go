package g

import (
	"encoding/json"
	"fmt"
	"github.com/mglslg/gpt-play/cmd/g/ds"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"time"
)

var Logger *log.Logger
var Conf ds.GlobalConfig
var SecToken ds.Token
var Role ds.Role
var PrivateChatAuth ds.PrivateChatAuth

// InitConfig readConfig reads the config file and unmarshals it into the config variable
func InitConfig(configPath string) {
	fmt.Println("Reading config file...")

	file, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Read config failed...", err)
		return
	}

	err = yaml.Unmarshal(file, &Conf)

	if err != nil {
		fmt.Println("Resolve config file failed!", err)
		return
	}

	fmt.Println("Config file read successfully!")
}

func InitRole(roleName string) {
	roleConfFile := fmt.Sprintf("role/%s.json", roleName)

	file, err := os.ReadFile(roleConfFile)
	if err != nil {
		fmt.Println("Read role config failed:", err)
	}

	Role.Name = roleName
	err = json.Unmarshal(file, &Role)

	if err != nil {
		fmt.Println("Resolve role config file failed:", err)
	}
	fmt.Println("This is " + Role.Name)
}

func InitLogger() *os.File {
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("%s/deploy/logs/%s.log", Conf.Home, currentDate+"-"+Role.Name)

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

	file, err := os.ReadFile("config/role_secrets/" + Role.Name + ".yaml")

	if err != nil {
		Logger.Fatal(err.Error())
	}

	err = yaml.Unmarshal(file, &SecToken)

	if err != nil {
		Logger.Fatal(err.Error())
	}

	Logger.Println("Secret Config file read successfully!Token:", SecToken.Discord)
}

func InitPrivateChatAuth() {
	fmt.Println("Reading private chat authorize file...")
	file, err := os.ReadFile("config/authorize/private_chat.json")
	if err != nil {
		Logger.Fatal(err.Error())
	}
	err = json.Unmarshal(file, &PrivateChatAuth)
	if err != nil {
		Logger.Fatal(err.Error())
	}
	Logger.Println("private chat authorize read successfully!")
}
