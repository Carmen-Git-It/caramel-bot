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
			fmt.Println("Error querying the professor given, please try another professor name")
			fmt.Println(err)
			var message = "**ERROR**: Professor not found, please try again."
			err = dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
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
				fmt.Println("Error querying the professor given, please try another professor name")
				fmt.Println(err)
				var message = "**ERROR**: Professor not found, please try again."
				err = dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
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
				Description: rmp1.professorName + " vs " + rmp2.professorName,
				Color:       0x00ff00,
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