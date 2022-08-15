package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	compliment := Command{
		Name: "compliment",
		Help: "Used to compliment somebody",
		Exec: compliment,
	}
	commands["compliment"] = compliment
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

	// If owner is mentioned, insult the sender
	if len(m.Mentions) > 0 && m.Mentions[0].ID == "246732655373189120" {
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> bitch")
		return
	}

	// If no one was tagged, send a general "bitch" out into the world
	if len(messageList) == 1 {
		s.ChannelMessageSend(m.ChannelID, "bitch")
	} else if len(m.Mentions) > 0 {
		// else if someone was tagged, tag that person and call them a "bitch"
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Mentions[0].ID+"> bitch")
	}
}

func compliment(s *discordgo.Session, m *discordgo.MessageCreate, messageList []string) {
	const URI = "https://complimentr.com/api"

	if len(m.Mentions) == 0 {
		return
	}

	response, err := http.Get(URI)
	if err != nil {
		fmt.Println("Error contacting Compliment API server")
		return
	}

	defer response.Body.Close()

	var body Compliment

	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		fmt.Println("Error decoding JSON")
		return
	}

	if len(m.Mentions) > 0 {
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Mentions[0].ID+"> "+body.Compliment)
	}
}
