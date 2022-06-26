package main

import (
	"log"

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
		dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This is the first tester function, good job!",
			},
		})
	},
	"test-command2": func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
		dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This is the second tester function, yay!",
			},
		})
	}}

func addHandlers(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	dg.AddHandler(func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(dg, i)
		}
	})
}

func registerCommands(dg *discordgo.Session) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		command, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Println("Error! Cannot create command!")
		}
		registeredCommands[i] = command
	}
}
