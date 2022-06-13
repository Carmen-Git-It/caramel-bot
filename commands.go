package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = make(map[string]Command)
)

type Command struct {
	Name string
	Help string

	Exec func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

func init() {
	bitch := Command{
		Name: "bitch",
		Help: "Used to call someone a bitch",
		Exec: bitch,
	}
	commands["bitch"] = bitch
}

func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate, message string) {

	messageList := strings.Fields(message)

	// Discard if there are no nouns or verbs after the prefix
	if len(messageList) == 0 {
		return
	}

	var commandName = strings.ToLower(messageList[0]) // Grab the name of the command

	command, ok := commands[commandName]

	if ok {
		command.Exec(s, m, messageList)
	} else {
		fmt.Println("Error executing command")
	}
}

// Command to call someone a bitch
func bitch(s *discordgo.Session, m *discordgo.MessageCreate, messageList []string) {
	if len(m.Mentions) == 0 {
		return
	}

	if m.Mentions[0].ID == "246732655373189120" {
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> bitch")
		return
	}

	fmt.Println(messageList)
	fmt.Println(m.Mentions[0].ID)
	if len(messageList) == 1 {
		s.ChannelMessageSend(m.ChannelID, "bitch")
	} else {
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Mentions[0].ID+"> bitch")
	}
}
