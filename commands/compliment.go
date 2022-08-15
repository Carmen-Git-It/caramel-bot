package commands

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type Compliment struct {
	Compliment string
}

func CommandCompliment(dg *discordgo.Session, i *discordgo.InteractionCreate) {
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
}
