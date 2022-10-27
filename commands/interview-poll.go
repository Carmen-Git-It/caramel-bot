package commands

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	badger "github.com/dgraph-io/badger/v3"
)

type State struct {
	Votes struct {
		One     discordgo.User `json:"1"`
		Two     discordgo.User `json:"2"`
		Three   discordgo.User `json:"3"`
		Four    discordgo.User `json:"4"`
		Five    discordgo.User `json:"5"`
		Six     discordgo.User `json:"6"`
		Seven   discordgo.User `json:"7"`
		Eight   discordgo.User `json:"8"`
		Nine    discordgo.User `json:"9"`
		TenPlus discordgo.User `json:"10+"`
	}
}

const DBLabel = "interviewVotes"

// interview container type
type VotesContainer = map[int][]discordgo.User

func CommandInterview(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	db := openDB()
	defer db.Close()

	votes := getVotesFromDB(db)
	// fmt.Printf("CommandInterview votes: %v\n", votes) // __AUTO_GENERATED_PRINT_VAR__

	options := ParseUserOptions(dg, i)
	vote := options["vote"].IntValue()
	votes = removeUserVotes(votes, *i.Member.User)
	votes = addVote(votes, int(vote), *i.Member.User)
	fmt.Printf("CommandInterview votes: %+v\n", votes) // __AUTO_GENERATED_PRINT_VAR__
	err := saveVote(db, votes)
	if err == badger.ErrConflict {
		panic(err)
	}

	dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "hello",
		},
	})
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
func removeUserVotes(votes VotesContainer, user discordgo.User) VotesContainer {
	filteredContainer := make(VotesContainer)
	for voteCount, vote := range votes {
		for _, u := range vote {
			// add the users that are not the user we are removing
			if u.ID != user.ID {
				filteredContainer[voteCount] = append(filteredContainer[voteCount], u)
			}
		}
	}
	return filteredContainer
}

func removeIndex(s []discordgo.User, index int) []discordgo.User {
	return append(s[:index], s[index+1:]...)
}

// saveVote save the db to local disk
// db     : the db to save
// votes  : the votes to save
// returns: error if any
func saveVote(db *badger.DB, votes VotesContainer) error {
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
