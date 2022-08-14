package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const prefix = "!"

// Variables used for command line params
var (
	Token string
)

// Parse arguments
func init() {
	// Accept the Discord bot token from the command line
	// flag.StringVar(&Token, "t", "", "Bot Token")
	// flag.Parse()

	// Load the .env file
	err := godotenv.Load("token.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
	Token = os.Getenv("TOKEN")
}

// Handles any message being created in the guild, parses them,
// and sends them to the commands module.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Exit function if the message was created by a bot
	if m.Author.Bot {
		return
	}

	// Do nothing if prefix is not present
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// Parse the command by trimming the prefix
	parseCommand(s, m, strings.TrimPrefix(m.Content, prefix))

	// Log some details
	fmt.Println("Message received!\nAuthor: " + m.Author.Username + "\nMessage: " + m.Message.Content)
}

func main() {

	// Create a new discord session
	dg, err := discordgo.New("Bot " + Token)
	dg.AddHandler(addHandlers)
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

	registerCommands(dg)

	defer dg.Close()

	// Listen until signal is received to end.
	fmt.Println("Caramel Bot is running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)

	//Testing logic, remove later
	QueryProfessor("Fardad Soleimanloo")

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	removeCommands(dg)
}
