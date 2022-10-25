package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CommandInterview(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	fmt.Printf("CommandInterview options: %+v\n", options[0]) // __AUTO_GENERATED_PRINT_VAR__
	dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "hello",
		},
	})
}
