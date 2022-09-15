package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"time"
)

// lists games in the games list and veto lists
func ListGames() string {
	//this returns a string with the game picks and vetos neatly ordered
	message := "Current Pool:\n"
	for i, s := range gl {
		if !inslice(s, vl) {
			message = message + "     " + strconv.Itoa(i+1) + ". " + gl[i] + "\n"
		}
	}
	message = message + "Current vetos:\n"
	for i := 0; i < len(vl); i++ {
		message = message + "     " + strconv.Itoa(i+1) + ". " + vl[i] + "\n"
	}
	return message
}

// clears games list and veto list to reset everything
func Iclear() (err error) {
	gl = nil
	vl = nil
	//s.ChannelMessageDelete(m.ChannelID, m.ID)
	return err
}

// this adds a game to the veto list or takes the number shown by listgames, and adds the game at that point to the veto list.
// games on the veto list cannot be selected when rolling
func IVeto(Content string) (message string, err error) {
	veto := string(Content)[6:]
	intVeto, Verr := strconv.Atoi(veto)
	match := false
	//if veto input was a number then set the veto to the game in gl
	if Verr == nil {
		if intVeto-1 < len(gl) {
			veto = gl[intVeto-1] //intVeto given by the user is indexed at 1 not zero
			match = true
		}

	} else {
		for i := 0; i < len(gl); i++ {
			if veto == gl[i] {
				match = true
			}
		}
	}
	//if user messed up the bot command and the veto couldnt be added, call them out.
	if match {
		vl = append(vl, veto)
		message = veto + " vetoed\n" + ListGames()
	} else {
		message = "No match for veto found, try again idiot"
	}

	return message, err
}

// this can be called to remove a member from a string array, if a member was mistakenly added.
func Rm(Content string, slice []string) (NewSlice []string, err error) {

	index, Verr := strconv.Atoi(string(Content)[5:])
	index = index - 1 //make index start at 0
	if Verr == nil && len(slice) > 0 && index >= 0 {
		if index == 0 {
			NewSlice = slice[1:]
		} else if index < len(slice) {
			NewSlice = append(slice[:index], slice[index+1:]...)
		} else if index == len(slice) {
			NewSlice = slice[:len(slice)-1]
		}

	}
	return NewSlice, err
}

// this checks the presense for string n in string array h
func inslice(n string, h []string) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}

// this adds a pick to the games list
func IPick(Content string) (message string, err error) {
	pick := string(Content)[6:]
	gl = append(gl, pick)
	message = pick + " added\n" + ListGames()

	return message, err
}

// this function gathers the games in the games list, checks that they are not present in the veto list.
// Then builds a hash input from the current date and the names of all the games not vetoed.
// then generates a sudo random number from the input, uses the modulo of the number of games to select a game.
func IRoll() (message string, err error) {
	selections := make([]string, 0)

	//checks lack of presense in veto list
	for _, s := range gl {
		if !inslice(s, vl) {
			selections = append(selections, s)
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
	randomN, err := strconv.ParseInt(sha1_hash, 16, 64)
	//with get the remainder of the sudo random number by number of games.
	intPick := randomN % int64(len(selections))
	//define picked game as the index
	pickedGame := selections[intPick]
	//message out
	message = pickedGame + " has been decided out of:\n" + ListGames()
	return message, err
}
