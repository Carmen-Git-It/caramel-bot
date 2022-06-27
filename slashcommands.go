package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

const testGuildID = "985707181854826497"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "test-command",
		Description: "A simple test command",
	},
	{
		Name:        "test-command2",
		Description: "A second simple test command",
	},
	{
		Name:        "compliment-user",
		Description: "Give another user a compliment",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The user you would like to compliment",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
		},
	},
}

var CommandHandlers = map[string]func(dg *discordgo.Session, i *discordgo.InteractionCreate){
	"test-command": func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
		err := dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This is the first tester function, good job!",
			},
		})
		if err != nil {
			fmt.Println("Error with function test-command:")
			fmt.Println(err)
		}
	},
	"test-command2": func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
		err := dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This is the second tester function, yay!",
			},
		})
		if err != nil {
			fmt.Println("Error with function test-command2:")
			fmt.Println(err)
		}
	},
	"compliment-user": func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
		const URI = "https://complimentr.com/api"

		// Access API server to attempt to retrieve a compliment
		response, err := http.Get(URI)
		if err != nil {
			fmt.Println("Error contacting Compliment API server")
			fmt.Println(err)
			return
		}

		// Close the response after function end
		defer response.Body.Close()

		// Compliment struct to hold compliment information
		var compliment Compliment

		// Decode the JSON response into a Compliment struct
		err = json.NewDecoder(response.Body).Decode(&compliment)
		if err != nil {
			fmt.Println("Error decoding JSON")
			return
		}

		// Get options from the application data
		options := i.ApplicationCommandData().Options
		var message = ""

		// If the option exists, add the result of the user option to the message as a ping
		if options[0] != nil {
			message = "<@" + options[0].UserValue(nil).ID + "> " + compliment.Compliment
		}

		// Build out the interaction response
		err = dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})

		if err != nil {
			fmt.Println("Error with function compliment:")
			fmt.Println(err)
		}
	}}

func addHandlers(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	dg.AddHandler(func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(dg, i)
		} else {
			fmt.Println("Error adding handler")
		}
	})
}

func registerCommands(dg *discordgo.Session) {
	fmt.Println("Registering commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for _, g := range dg.State.Guilds {
		for i, v := range Commands {
			command, err := dg.ApplicationCommandCreate(dg.State.User.ID, g.ID, v)
			if err != nil {
				fmt.Println("Error! Cannot create command!")
				fmt.Println(err)
			}
			registeredCommands[i] = command
		}
	}

}
