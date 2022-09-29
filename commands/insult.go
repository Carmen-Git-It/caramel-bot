package commands

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type Insult struct {
	Insult string `json:"insult"`
}

func CommandInsult(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	const URI = "https://insult.mattbas.org/api/insult"

	// Send get request
	response, err := http.Get(URI)
	if err != nil {
		fmt.Println("Error contacting Insult API server")
		fmt.Println(err)
		return
	}

	defer response.Body.Close()
	var insult Insult

	err = json.NewDecoder(response.Body).Decode(&insult)
	if err != nil {
		fmt.Println("Error decoding JSON")
		return
	}

	options := i.ApplicationCommandData().Options
	var message = ""

	// If the option exists
	if options[0] != nil {
		message = "<@" + options[0].UserValue(nil).ID + "> " + insult.Insult
	}

	err = dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		fmt.Println("Error with function insult:")
		fmt.Println(err)
	}
}
