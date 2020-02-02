package bot

import (
	"encoding/json"
	"fmt"
	"github.com/DVC-Software/discordbot/execution"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
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

// return the eparsed command and remove that command from the message
// Return tha matched command and the striped message
func ParseCommand(message string) (string, string) {
	prefix := regexp.MustCompile(`^/(.*)$`)
	// every message should have the prefix now
	if prefix.MatchString(message) {
		fmt.Println("command prefix matched")
		message = prefix.ReplaceAllString(message, "$1")
	}
	// check create member
	hello := regexp.MustCompile(`^hello$`)
	if hello.MatchString(message) {
		fmt.Println("hello matched")
		return "hello", ""
	}
	// check create member
	create_member := regexp.MustCompile(`^create member[\s]*(.*)$`)
	if create_member.MatchString(message) {
		fmt.Println("create member matched")
		message = create_member.ReplaceAllString(message, "$1")
		return "create_member", message
	}

	return "invalid", message
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
	command, message := ParseCommand(m.Content)
	switch command {
	case "create_member":
		args := []string{message, m.Author.ID}
		msg, err := execution.CreateMember(args)
		_, _ = s.ChannelMessageSend(m.ChannelID, msg+"\n"+err)
	case "hello":
		valid, info := execution.IdentifyMember(m.Author.ID)
		var msg string
		if valid {
			msg = info.Name
		} else {
			msg = m.Author.Username
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, "Hello, "+msg)
	default:
		_, _ = s.ChannelMessageSend(m.ChannelID, "Oops, this is not a valid command!")
	}
}
