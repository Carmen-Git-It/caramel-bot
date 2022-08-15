package main

import (
	"fmt"

	c "caramel-bot/commands"

	"github.com/bwmarrin/discordgo"
)

var registeredCommands []*discordgo.ApplicationCommand

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "compliment",
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
	{
		Name:        "bitch",
		Description: "Call another user a bitch",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The user you would like to call a bitch",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
		},
	},
	{
		Name:        "rmp",
		Description: "Query RateMyProfessors for data on a professor",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "professor",
				Description: "The professor you would like to look up",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
	{
		Name:        "rmp-compare",
		Description: "Compare two professors using data from RateMyProfessors",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "first-professor",
				Description: "The first professor you would like to compare",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "second-professor",
				Description: "The second professor that you would like to compare",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
}

var CommandHandlers = map[string]func(dg *discordgo.Session, i *discordgo.InteractionCreate){
	"compliment":  c.CommandCompliment,
	"bitch":       c.CommandBitch,
	"rmp":         c.CommandRMP,
	"rmp-compare": c.CommandRMPCompare,
}

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
	registeredCommands = make([]*discordgo.ApplicationCommand, len(Commands))
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

func removeCommands(dg *discordgo.Session) {
	for _, g := range dg.State.Guilds {
		for _, v := range registeredCommands {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, g.ID, v.ID)
			if err != nil {
				fmt.Println("Error! Cannot delete command!")
				fmt.Println(err)
			}
		}
	}
}
