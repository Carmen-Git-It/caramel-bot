package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CommandRMPCompare(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// Query both professors, return if error encountered
	if options[0] != nil && options[1] != nil {
		rmp1, err := QueryProfessor(options[0].StringValue())
		if err != nil {
			fmt.Println("Error querying professor" + options[0].StringValue())
			fmt.Println(err)
			var message = "Could not find professor \"" + options[0].StringValue() + "\", please try again."
			// Make errors visible only to the one using the command
			err = dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   1 << 6, // Ephemeral
					Content: message,
				},
			})

			if err != nil {
				fmt.Println("Error with function rmp-compare")
				fmt.Println(err)
				return
			}
		} else {
			rmp2, err := QueryProfessor(options[1].StringValue())
			if err != nil {
				fmt.Println("Error querying professor" + options[1].StringValue())
				fmt.Println(err)
				var message = "Could not find professor \"" + options[1].StringValue() + "\", please try again."
				err = dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   1 << 6, // Ephemeral
						Content: message,
					},
				})

				if err != nil {
					fmt.Println("Error with function rmp-compare")
					fmt.Println(err)
					return
				}
			}

			// Have both professors data, now compare them
			embed := &discordgo.MessageEmbed{
				Title:       "RateMyProfessor Comparison",
				Description: "[" + rmp1.professorName + "](" + rmp1.rmpURL + ")" + " vs " + "[" + rmp2.professorName + "](" + rmp2.rmpURL + ")",
				Color:       0x00ff00,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Data retreived from RateMyProfessors.com",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Best Overall Rating",
						Value:  CompareOverallRating(rmp1, rmp2),
						Inline: false,
					},
					{
						Name:   "Highest Would Take Again %",
						Value:  CompareWouldTakeAgain(rmp1, rmp2),
						Inline: false,
					},
					{
						Name:   "Lowest Difficulty",
						Value:  CompareDifficulty(rmp1, rmp2),
						Inline: false,
					},
					{
						Name:   "Best Rating by Course",
						Value:  CompareBestByCourse(rmp1, rmp2),
						Inline: false,
					},
				},
			}
			err = dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds:  []*discordgo.MessageEmbed{embed},
					Content: "",
				},
			})

			if err != nil {
				fmt.Println("Error with function rmp-compare")
				fmt.Println(err)
			}
		}
	}
}
