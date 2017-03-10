package amagi

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	// "github.com/b-eee/amagi"
	"github.com/b-eee/amagi/api/pubnub"
	"github.com/b-eee/amagi/api/slack"
	"github.com/joho/godotenv"
)

var (
	configEnvFileName = "config.env"
)

// InitLogCredentials initialize logger credentials
func InitLogCredentials() slack.Host {
	if err := initializeConfigEnv(); err != nil {
		panic("please set your config.env!")
	}

	return slack.Host{
		TokenID:   os.Getenv("SLACK_TOKEN"),
		ChannelID: os.Getenv("SLACK_CHANNEL_ID"),
		Hostname: func() string {
			return pubnub.GetCurrentHostIP()
		},
		Env: os.Getenv("ENV"),

		PublishKey:   os.Getenv("PUBLISH_KEY"),
		SubscribeKey: os.Getenv("SUBSCRIBE_KEY"),
		SecretKey:    os.Getenv("SECRET_KEY"),
		MicroAppName: os.Getenv("APP_NAME"),
	}
}

func initializeConfigEnv() error {
	file := fmt.Sprintf("%v/%v", getCwd(), configEnvFileName)
	if err := godotenv.Load(file); err != nil {
		log.Printf("ERROR LOADING DOTENV %v at path %v", err, file)
	}

	fmt.Println("config.env initialize!")
	return nil
}

func getCwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// if path is root, return empty instead
	if pwd == "/" {
		pwd = ""
	}

	return pwd
}
