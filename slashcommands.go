package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "test-command",
		Description: "A simple test command",
	},
	{
		Name:        "test-command2",
		Description: "A second simple test command",
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
	for i, v := range Commands {
		command, err := dg.ApplicationCommandCreate(dg.State.User.ID, "985707181854826497", v)
		if err != nil {
			fmt.Println("Error! Cannot create command!")
			fmt.Println(err)
		}
		registeredCommands[i] = command
	}

}
