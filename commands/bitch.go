package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CommandBitch(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	var message = ""

	// options := i.ApplicationCommandData().Options
	options := ParseUserOptions(dg, i)
	user := options["user"]

	if user != nil && user.UserValue(dg).ID != "246732655373189120" {
		message = "<@" + user.UserValue(dg).ID + "> is a bitch."
	} else if user != nil {
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
