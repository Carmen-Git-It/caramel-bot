package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CommandBitch(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	var message = ""

	options := i.ApplicationCommandData().Options

	if options[0] != nil && options[0].UserValue(nil).ID != "246732655373189120" {
		message = "<@" + options[0].UserValue(nil).ID + "> is a bitch."
	} else if options[0] != nil {
		message = "<@" + i.Member.User.ID + "> nice try, you're a bitch."
	}

	err := dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		fmt.Println("Error with function bitch:")
		fmt.Println(err)
	}
}
