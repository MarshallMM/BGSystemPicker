package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token    string
	gameList []Game
)

type Game struct {
	name     string
	veto     bool
	pickedBy string
	vetoedBy string
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}
func main() {
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
		//prevPost = m.MessageReference.MessageID
		//prevChan = m.ChannelID
		return
	}

	if mes == "!clear" {
		Iclear()
	}
	if mes == "!list" {
		message = ListGames()
	}

	if mes == "!trout" {
		message = "trout that"
	}
	if mes == "!roll" {
		fmt.Println("rolling")
		message, err = IRoll()
	}

	if len(mes) > 5 {
		if mes[:4] == "!rmp" {
			message = Rmp(mes)
		}
		if mes[:4] == "!rmv" {
			message = Rmv(mes)

		}
		if mes[:5] == "!pick" {
			message, err = IPick(mes)
		}

		if mes[:5] == "!veto" {
			message, err = IVeto(mes)
		}
	}

	if message != "hmm" {
		//err = s.ChannelMessageDelete(prevChan, prevPost)
		if err != nil {
			fmt.Println(err)
		}
		// Send a text message
		_, err = s.ChannelMessageSend(m.ChannelID, message)

		if err != nil {
			fmt.Println(err)
		}

	} else {
		return
	}
}
