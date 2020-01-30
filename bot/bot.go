package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

var BotID string
var bot *discordgo.Session

func GetBot() *discordgo.Session {
	return bot
}

func GetBotID() string {
	return BotID
}

func Start() {
	bot, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		panic(err.Error())
	}
	u, err := bot.User("@me")
	if err != nil {
		panic(err.Error())
	}
	// configure settings
	bot.ShouldReconnectOnError = true
	bot.LogLevel = 1

	// Add handlers
	bot.AddHandler(MessageHandler)

	BotID = u.ID
	err = bot.Open()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Bot online!" + BotID)
}
