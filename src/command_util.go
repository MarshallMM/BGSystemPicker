package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// lists games in the games list and veto lists
func ListGames() string {
	//this returns a string with the game picks and vetos neatly ordered
	message := "Current Pool:\n"
	vetoMessage := ""
	//sort non vetoed games first
	sort.Slice(gameList[:], func(i, j int) bool {
		return gameList[i].veto < gameList[j].veto
	})
	for i, s := range gameList {
		if s.veto == 1 {
			vetoMessage = " [vetoed by " + s.vetoedBy + "]"
		} else {
			vetoMessage = ""
		}
		message = message + "     " + strconv.Itoa(i+1) + ". " + s.name + " [picked by " + s.pickedBy + "]"
		message = message + vetoMessage + "\n"
	}

	return message
}

// clears games list and veto list to reset everything
func Iclear() (err error) {
	gameList = nil
	//s.ChannelMessageDelete(m.ChannelID, m.ID)
	return err
}

// this adds a pick to the games list
func IPick(Content string, m *discordgo.MessageCreate) (message string) {
	pick := string(Content)[6:]
	gameList = append(gameList, Game{
		name:     pick,
		veto:     0,
		pickedBy: m.Author.Username,
		vetoedBy: "",
	})

	message = pick + " added\n" + ListGames()

	return message
}

// this adds a game to the veto list or takes the number shown by listgames, and adds the game at that point to the veto list.
// games on the veto list cannot be selected when rolling
func IVeto(Content string, m *discordgo.MessageCreate) (message string) {
	veto := string(Content)[6:]
	intVeto, vErr := strconv.Atoi(veto)

	if vErr == nil {

	} else { //turn a game input into a number
		for i, s := range gameList {
			if veto == s.name {
				intVeto = i
			}
		}
	}
	intVeto = intVeto - 1 //intVeto given by the user is indexed at 1 not zero
	if intVeto < 0 || intVeto >= len(gameList) {
		if gameList[intVeto].veto == 0 {
			gameList[intVeto].veto = 1
			gameList[intVeto].vetoedBy = m.Author.Username
			return veto + " vetoed\n" + ListGames()
		} else {
			return "already vetoed numnuts"
		}
	}
	//if user messed up the bot command and the veto couldnt be added, call them out.
	return "No match for veto found, try again idiot"
}

// this can be called to remove a member from an array, if a member was mistakenly added.
func Rmp(Content string) (message string) {

	index, Verr := strconv.Atoi(string(Content)[5:])
	index = index - 1 //make index start at 0, user inputs start at 1
	if Verr == nil && len(gameList) > 0 && index >= 0 {
		if index == 0 {
			gameList = gameList[1:]
			message = ListGames()
		} else if index < len(gameList) {
			gameList = append(gameList[:index], gameList[index+1:]...)
			message = ListGames()
		} else if index == len(gameList) {
			gameList = gameList[:len(gameList)-1]
			message = ListGames()
		} else {
			message = "idk what number that was but it dont work"
		}

	} else {
		message = "Not an integer input"
	}
	return message
}
func Rmv(Content string) (message string) {
	index, Verr := strconv.Atoi(string(Content)[5:])
	index = index - 1 //make index start at 0, user inputs start at 1
	if Verr == nil {
		if gameList[index].veto == 1 {
			gameList[index].veto = 0
			gameList[index].vetoedBy = ""
			message = ListGames()
		} else {
			message = "idk what number that was but it dont work"
		}
	} else {
		message = "Not an integer input"
	}
	return message
}

// this function gathers the games in the games list, checks that they are not present in the veto list.
// Then builds a hash input from the current date and the names of all the games not vetoed.
// then generates a sudo random number from the input, uses the modulo of the number of games to select a game.
func IRoll() (message string) {
	selections := make([]string, 0)

	//checks lack of presense in veto list
	for _, s := range gameList {
		if s.veto == 0 {
			selections = append(selections, s.name)
		}
	}
	//sorts selections so to avoid order of picks effecting result.
	sort.Strings(selections)
	fmt.Println(selections)
	//create hash input with date then add selections
	hash := time.Now().Format("01-02-2006")
	for i := 0; i < len(selections); i++ {
		hash = hash + selections[i]
	}
	fmt.Println(hash)

	//get a sudo random number from the input
	h := sha1.New()
	h.Write([]byte(hash))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	fmt.Println(hash, sha1_hash)
	randomN, _ := strconv.ParseInt(sha1_hash, 16, 64)
	//with get the remainder of the sudo random number by number of games.
	intPick := randomN % int64(len(selections))
	//define picked game as the index
	pickedGame := selections[intPick]
	//message out
	message = pickedGame + " has been decided out of:\n" + ListGames()
	return message
}
