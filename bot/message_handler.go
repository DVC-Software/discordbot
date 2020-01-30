package bot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func NormalizeMessage(message string) string {
	// filter out prefixing @ mentions before the ! command prefix
	fmt.Println(message)
	front_mention := regexp.MustCompile(`^<@!*(.*?)/(.*)$`)
	rear_mention := regexp.MustCompile(`/(.*)<@!*(.*?)$`)
	multiple_prefixes := regexp.MustCompile(`(.*)//+(.*)`)
	if front_mention.MatchString(message) {
		fmt.Println("front mention matched")
		message = front_mention.ReplaceAllString(message, "/$2")
	} else if rear_mention.MatchString(message) {
		fmt.Println("rearmention matched")
		message = rear_mention.ReplaceAllString(message, "/$1")
	}
	if multiple_prefixes.MatchString(message) {
		fmt.Println("multiple prefixes matched")
		message = ""
	}
	// to lower case
	message = strings.TrimSpace(strings.ToLower(message))
	return message
}

func ValidateMessage(m *discordgo.MessageCreate) (bool, string) {
	if m.Author.ID == GetBotID() {
		return false, ""
	}
	// check if the message is a dm
	if m.GuildID != "" {
		// check if the message mentions the bot
		mentioned := false
		for _, usr := range m.Mentions {
			if usr.ID == GetBotID() {
				mentioned = true
				break
			}
		}
		if !mentioned {
			return false, ""
		}
	}
	if !strings.HasPrefix(m.Content, "/") || m.Content == "/" {
		return false, "Please use command staring with a '/'"
	}
	return true, ""
}

func ParseCommand(message string) string {
	prefix := regexp.MustCompile(`^/(.*)$`)
	if prefix.MatchString(message) {
		fmt.Println("command prefix matched")
		message = prefix.ReplaceAllString(message, "$1")
	}
	return message
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = NormalizeMessage(m.Content)
	valid, errMsg := ValidateMessage(m)
	if !valid {
		// write error message
		if errMsg != "" {
			_, _ = s.ChannelMessageSend(m.ChannelID, errMsg)
		}
		return
	}
	// print the message json
	msgJson, _ := json.Marshal(m)
	fmt.Println(string(msgJson))
	// use regex to distinguish different commands
	switch command := ParseCommand(m.Content); command {
	case "create member":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Create member command is not implemented now!")
	case "hello":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Hello, "+m.Author.Username)
	default:
		_, _ = s.ChannelMessageSend(m.ChannelID, "Oops, this is not a valid command!")
	}
}
