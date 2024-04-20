package g

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
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
var SessionMap map[string]*ds.UserSession

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
	logPath := fmt.Sprintf("%s/logs", Conf.Home)

	// 检查logs目录是否存在，如果不存在则创建
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(logPath, 0755); mkErr != nil {
			log.Fatalf("Unable to create log directory: %v", mkErr)
		}
	}

	logFileName := fmt.Sprintf("%s/%s-%s.log", logPath, currentDate, Role.Name)

	// 创建一个日志文件
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Unable to open log file: %v", err)
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

func InitSessionMap() {
	SessionMap = make(map[string]*ds.UserSession)
}

// GetUserSession 获取当前用户session,如果没有则创建
func GetUserSession(authorId string, authorName string, channelId string, s *discordgo.Session) *ds.UserSession {
	channel, err := s.Channel(channelId)
	if err != nil {
		Logger.Fatal("获取频道信息失败", err)
	}

	key := getUserChannelId(authorId, channelId)
	_, exists := SessionMap[key]
	if !exists {
		SessionMap[key] = newUserSession(authorId, authorName, channelId, channel.ParentID)
	}
	return SessionMap[key]
}

func newUserSession(authorId string, authorName string, channelId string, parentChannelId string) *ds.UserSession {
	userChannelId := getUserChannelId(authorId, channelId)
	return &ds.UserSession{
		UserId:          authorId,
		UserName:        authorName,
		UserChannelID:   userChannelId,
		ChannelID:       channelId,
		ParentChannelID: parentChannelId,
		ClearDelimiter:  Role.ClearDelimiter,
		Model:           "gpt-3.5-turbo-0125",
		Temperature:     0.7,
		Prompt:          Role.Characters[0].Desc,
		AllowChannelIds: Role.ChannelIds,
		OnConversation:  true, //默认都是有上下文的
		OnAt:            true, //默认都是需要艾特才回复的
	}
}

func getUserChannelId(authorId string, channelId string) string {
	return authorId + "_" + channelId
}
