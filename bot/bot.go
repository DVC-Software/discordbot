package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

var BotID string
var bot *discordgo.Session

func GetBot() (*discordgo.Session, string) {
	return bot, BotID
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

	BotID = u.ID
	err = bot.Open()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Bot online!" + BotID)
}
