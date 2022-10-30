package commands

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/olekukonko/tablewriter"
)

// the label of the db
const DBLabel = "interviewVotes"

// interview container type
type VotesContainer = map[int][]discordgo.User

func CommandInterview(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	db := openDB()
	defer db.Close()

	options := ParseUserOptions(dg, i)
	// record a user vote
	if val, ok := options["vote"]; ok {
		vote := val.IntValue()
		votes := getVotesFromDB(db)

		mess := formatVotes(votes, "Recording your vote... It may take a few seconds")

		dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: mess,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		votes, _ = removeUserVotes(votes, *i.Member.User)
		votes = addVote(votes, int(vote), *i.Member.User)
		err := saveVotes(db, votes)
		if err == badger.ErrConflict {
			panic(err)
		}

		time.Sleep(3 * time.Second)

		mess = formatVotes(votes, "Your vote has been recorded successfully")
		dg.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &mess,
		})

	} else if val, ok := options["getvotes"]; ok {
		// return all user votes
		if val.BoolValue() {
			votes := getVotesFromDB(db)
			mess := formatVotes(votes, "")

			dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: mess,
				},
			})
		} else {
			dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ok",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

		}
	} else if val, ok := options["remove"]; ok {
		// remove a user's votes
		if val.BoolValue() {
			votes := getVotesFromDB(db)
			message := formatVotes(votes, "Removing your vote... This may take a few seconds")
			dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			votes, removed := removeUserVotes(votes, *i.Member.User)
			saveVotes(db, votes)

			time.Sleep(3 * time.Second)

			// if a user's vote was removed, tell them their vote has been removed
			if removed {
				message = formatVotes(votes, "Your vote has been removed")
				dg.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &message,
				})
			} else {
				message = formatVotes(votes, "You have not voted... Pls vote")
				dg.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &message,
				})
			}
		} else {
			dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ok",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	}
}

// openDB opens a connection to the local database
func openDB() *badger.DB {
	opts := badger.DefaultOptions("./db")
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	return db
}

// getVotesFromDB retrieve the votes from local databse
// db: the db to retrieve from
func getVotesFromDB(db *badger.DB) VotesContainer {
	// container for votes
	votes := make(VotesContainer)
	db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(DBLabel))
		if err != nil {
			return err
		}
		// return nil
		err = item.Value(func(val []byte) error {
			err := json.Unmarshal(val, &votes)
			if err != nil {
				return err
			}
			return err
		})
		return err
	})

	// if err != badger.ErrKeyNotFound {
	// panic(err)
	// }
	return votes
}

// addVote add a new vote to the existing collection of votes, and return the votes container
// votes  : a map of existing votes
// newVote: the new vote to add, in numeric value
// user   : the user who voted
// returns: the updated votes container
func addVote(votes VotesContainer, newVote int, user discordgo.User) VotesContainer {
	votes[newVote] = append(votes[newVote], user)
	return votes
}

// removeUserVotes removes all of a user's votes
// votes  : the votes container
// user   : the user's votes to remove
// returns: the updated votes container
func removeUserVotes(votes VotesContainer, user discordgo.User) (VotesContainer, bool) {
	filteredContainer := make(VotesContainer)
	// if we encountered the user we are removing
	found := false
	for voteCount, vote := range votes {
		for _, u := range vote {
			// add the users that are not the user we are removing
			if u.ID != user.ID {
				filteredContainer[voteCount] = append(filteredContainer[voteCount], u)
			} else {
				// keep track of we've encountered the user we are removing
				found = true
			}
		}
	}
	return filteredContainer, found
}

func removeIndex(s []discordgo.User, index int) []discordgo.User {
	return append(s[:index], s[index+1:]...)
}

// saveVotes save the db to local disk
// db     : the db to save
// votes  : the votes to save
// returns: error if any
func saveVotes(db *badger.DB, votes VotesContainer) error {
	err := db.Update(func(txn *badger.Txn) error {
		j, err := json.Marshal(votes)
		if err != nil {
			return err
		}
		entry := badger.NewEntry([]byte(DBLabel), j)
		err = txn.SetEntry(entry)
		return err
	})
	return err
}

// formatVotes formats a container of votes as a discord message
// votes  : the votes to format
// message: a message to append to the end of the vote tally
// retturn: the formatted message
func formatVotes(votes VotesContainer, message string) string {
	data := make([][]string, 0)
	mess := &strings.Builder{}
	mess.WriteString("```")
	for idx, vote := range votes {
		// string representation of the index
		strIdx := strconv.Itoa(idx)
		// string representation of the number of votes
		strLen := strconv.Itoa(len(vote))
		data = append(data, []string{strIdx, strLen})
	}

	table := tablewriter.NewWriter(mess)
	table.SetHeader([]string{"NO.Interviews", "Votes"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	mess.WriteString("```")
	mess.WriteString("\n\n" + message)
	return mess.String()
}
