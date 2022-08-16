package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func CommandRMP(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	if options[0] != nil {
		rmp, err := QueryProfessor(options[0].StringValue())
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
				fmt.Println("Error with function rmp")
				fmt.Println(err)
			}
		} else {
			var ratingByCourse = ""
			var topTags = ""
			for _, tag := range rmp.topTags {
				topTags += tag + ", "
			}
			for _, course := range rmp.courses {
				ratingByCourse = fmt.Sprintf("%s**%s**%s%.2f%s%d%s", ratingByCourse, course, ": ", rmp.totalRatingByCourse[course], "/5 from ", rmp.numRatingsByCourse[course], " reviews\n")
			}

			var imageUrl string = fmt.Sprint("https://image-charts.com/chart?cht=bvg&chbr=10&chd=t:", rmp.ratingDistribution[1], ",", rmp.ratingDistribution[2], ",", rmp.ratingDistribution[3], ",", rmp.ratingDistribution[4], ",", rmp.ratingDistribution[5], "&chxr=0,1,6,1&chxt=x,y&chs=500x400&chdls=000000,18&chtt=Rating+Distrbution")

			// Testing code
			fmt.Println("Top tags: ", rmp.topTags)

			embed := &discordgo.MessageEmbed{
				Color:       0x00ff00,
				Description: "[View profile on RateMyProfessors.com](" + rmp.rmpURL + ")",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Overall rating",
						Value:  rmp.totalRating + "/5",
						Inline: true,
					},
					{
						Name:   "Would Take Again",
						Value:  rmp.wouldTakeAgain,
						Inline: true,
					},
					{
						Name:   "Difficulty",
						Value:  rmp.levelOfDifficulty + "/5",
						Inline: true,
					},
					{
						Name:   "Top Tags",
						Value:  topTags,
						Inline: false,
					},
					{
						Name:   "Average Rating by Course",
						Value:  ratingByCourse,
						Inline: false,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
				Title:     rmp.professorName,
				Image: &discordgo.MessageEmbedImage{
					URL: imageUrl,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Powered by image-charts.com | Data retreived from RateMyProfessors.com",
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
				fmt.Println("Error with function rmp")
				fmt.Println(err)
			}
		}
	}

}
