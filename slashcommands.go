package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	"compliment": func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
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
	},
	"bitch": func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
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
	},
	"rmp": func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
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
					ratingByCourse = fmt.Sprintf("%s%s%s%.2f%s%d%s", ratingByCourse, course, ": ", rmp.totalRatingByCourse[course], "/5 from ", rmp.numRatingsByCourse[course], " reviews\n")
				}

				var imageUrl string = fmt.Sprint("https://image-charts.com/chart?cht=bvg&chbr=10&chd=t:", rmp.ratingDistribution[1], ",", rmp.ratingDistribution[2], ",", rmp.ratingDistribution[3], ",", rmp.ratingDistribution[4], ",", rmp.ratingDistribution[5], "&chxr=0,1,6,1&chxt=x,y&chs=500x400&chdls=000000,18&chtt=Rating+Distrbution")

				// Testing code
				fmt.Println("Top tags: ", rmp.topTags)

				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Color:  0x00ff00,
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
						Text: "Powered by image-charts.com | Data retreived from RateMyProfessor.com",
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

	},
	"rmp-compare": func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
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

				// Compose the message in markdown for formatting purposes

				// Have both professors data, now compare them
				embed := &discordgo.MessageEmbed{
					Title:       "RateMyProfessor Comparison",
					Description: fmt.Sprintf("%s and %s", rmp1.professorName, rmp2.professorName),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "\u200B",
							Value:  "```**Name**         :**" + rmp1.professorName + "**     " + "**" + rmp2.professorName + "**```",
							Inline: false,
						},
						// // Line break
						// {
						// 	Name:   "\u200B",
						// 	Value:  "\u200B",
						// 	Inline: false,
						// },
						{
							Name:   "\u200B",
							Value:  "**Total Rating**:**" + rmp1.totalRating + "/5**     " + "**" + rmp2.totalRating + "/5**",
							Inline: true,
						},
						// Line break
						{
							Name:   "\u200B",
							Value:  "\u200B",
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
	},
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
