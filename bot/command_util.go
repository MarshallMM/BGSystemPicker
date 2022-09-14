package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"time"
)

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
func Iclear() (err error) {
	gl = nil
	vl = nil
	//s.ChannelMessageDelete(m.ChannelID, m.ID)
	return err
}

func IList() (message string, err error) {
	message = ListGames()
	return message, err
}
func IVeto(Content string) (message string, err error) {
	veto := string(Content)[6:]
	intVeto, Verr := strconv.Atoi(veto)
	match := false
	//if veto input was a number then set the veto to the game in gl
	if Verr == nil {
		if intVeto-1 < len(gl) {
			veto = gl[intVeto-1]
			match = true
		}

	} else {
		for i := 0; i < len(gl); i++ {
			if veto == gl[i] {
				match = true
			}
		}
	}
	if match {
		vl = append(vl, veto)
		message = veto + " vetoed\n" + ListGames()
	} else {
		message = "No match for veto found, try again idiot"
	}

	return message, err
}
func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
func Rm(Content string, slice []string) (NewSlice []string, err error) {

	index, Verr := strconv.Atoi(string(Content)[5:])
	index = index - 1 //make index start at 0
	if Verr == nil && len(slice) > 0 && index >= 0 {
		if index == 0 {
			NewSlice = slice[1:]
		} else if index < len(slice) {
			NewSlice = remove(slice, index)
		} else if index == len(slice) {
			NewSlice = slice[:len(slice)-1]
		}

	}
	return NewSlice, err
}

func inslice(n string, h []string) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}
func IPick(Content string) (message string, err error) {
	pick := string(Content)[6:]
	gl = append(gl, pick)
	message = pick + " added\n" + ListGames()

	return message, err
}

func IRoll() (message string, err error) {
	selections := make([]string, 0)

	for _, s := range gl {
		if !inslice(s, vl) {
			selections = append(selections, s)
		}
	}
	sort.Strings(selections)
	fmt.Println(selections)
	hash := time.Now().Format("01-02-2006")
	for i := 0; i < len(selections); i++ {
		hash = hash + selections[i]
	}
	fmt.Println(hash)

	h := sha1.New()
	h.Write([]byte(hash))
	sha1_hash := hex.EncodeToString(h.Sum(nil))

	fmt.Println(hash, sha1_hash)

	pick, err := strconv.ParseInt(sha1_hash, 16, 64)
	pick = pick % int64(len(selections))
	pickedGame := selections[pick]
	message = pickedGame
	message = pickedGame + " has been decided out of:\n" + ListGames()
	return message, err
}
