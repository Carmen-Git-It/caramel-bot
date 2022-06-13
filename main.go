package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line params
var (
	Token string
)

// Parse arguments
func init() {
	// Accept the Discord bot token from the command line
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	fmt.Println("Message received!\nAuthor: " + m.Author.Username + "\nMessage: " + m.Message.Content)
	_, err := s.ChannelMessageSend(m.ChannelID, "bitch")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	// Create a new discord session
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating new discord session, ", err)
		panic(err)
	}

	// Add a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Only cares about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening a connection to discord, ", err)
		panic(err)
	}

	// Listen until signal is received to end.
	fmt.Println("Caramel Bot is running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
