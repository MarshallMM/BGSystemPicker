package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
	gl    []string
	vl    []string
)

const KuteGoAPIURL = "https://kutego-api-xxxxx-ew.a.run.app"

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}
func main() {
	//var gameList []string
	//var vetoList []string
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type Gopher struct {
	Name string `json: "name"`
}

func listGames() string {
	//this returns a string with the game picks and vetos neatly ordered
	message := "Current Picks:\n"
	for i, s := range gl {
		if !inslice(s, vl) {
			message = message + "     " + strconv.Itoa(i+1) + ". " + gl[i] + "\n"
		}
	}
	message = message + "Current vetos:\n"
	for i := 0; i < len(vl); i++ {
		message = message + "     " + vl[i] + "\n"
	}
	return message
}
func iclear() (err error) {
	gl = nil
	vl = nil
	//s.ChannelMessageDelete(m.ChannelID, m.ID)
	return err
}

func iList() (message string, err error) {
	message = listGames()
	return message, err
}
func iVeto(Content string) (message string, err error) {
	veto := string(Content)[6:]
	intVeto, Verr := strconv.Atoi(veto)
	match := false
	//if veto input was a number then set the veto to the game in gl
	if Verr == nil {
		if intVeto-1 < len(gl) {
			veto = gl[intVeto-1]
			match = true
		} else {

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
		message = veto + " vetoed\n" + listGames()
	} else {
		message = "No match for veto found, try again idiot"
	}

	return message, err
}
func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
func rm(Content string) (message string, err error) {

	index, Verr := strconv.Atoi(string(Content)[4:])

	if Verr == nil && len(gl) > 0 && index >= 0 {
		if index < len(gl) {
			remove(gl, index)
		} else if index == len(gl) {
			gl = gl[:len(gl)-1]
		}
		message = "removed " + fmt.Sprint(index) + "\n" + listGames()
	}
	return message, err
}

func inslice(n string, h []string) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}
func iPick(Content string) (message string, err error) {
	pick := string(Content)[6:]
	gl = append(gl, pick)
	message = pick + " added\n" + listGames()

	return message, err
}

func iRoll() (message string, err error) {
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
	message = pickedGame + " has been decided out of:\n" + listGames()
	return message, err
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// This function will be called (due to AddHandler above) every time a new
	// message is created on any channel that the authenticated bot has access to.
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	message := "hmm"
	mes := strings.ToLower(m.Content)
	var err error
	err = nil
	if m.Author.ID == s.State.User.ID {
		return
	}

	if mes == "!clear" {
		iclear()
	}
	if mes == "!list" {
		message, err = iList()

	}

	if mes == "!trout" {
		message = "trout that"
	}
	if mes == "!roll" {
		fmt.Println("rolling")
		message, err = iRoll()
	}
	if len(mes) > 4 {
		if mes[:3] == "!rm" {
			message, err = rm(mes)
		}
	}
	if len(mes) > 5 {

		if mes[:5] == "!pick" {
			message, err = iPick(mes)
		}

		if mes[:5] == "!veto" {
			message, err = iVeto(mes)
		}
	}

	if message != "hmm" {
		// Send a text message
		_, err = s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		return
	}
	return
}
