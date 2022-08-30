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
			message = message + "     " + gl[i] + "\n"
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
	//if veto input was a number then set the veto to the game in gl
	if Verr == nil {
		if intVeto-1 < len(gl) {
			veto = gl[intVeto-1]
		} else {

		}

	}
	vl = append(vl, veto)
	message = veto + " vetoed\n" + listGames()
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
	var err error
	err = nil
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!clear" {
		iclear()
	}
	if m.Content == "!list" {
		message, err = iList()

	}
	if m.Content == "!trout" {
		message = "trout that"
	}
	if m.Content == "!roll" {
		fmt.Println("rolling")
		message, err = iRoll()
	}

	if len(m.Content) > 5 {
		if m.Content[:5] == "!pick" {
			message, err = iPick(m.Content)
		}

		if m.Content[:5] == "!veto" {
			message, err = iVeto(m.Content)
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
