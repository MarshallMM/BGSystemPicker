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
	Token             string
	gameList          []Game
	previousMessageID string
)

type Game struct {
	name     string
	veto     int
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
	var (
		err     error
		message string
	)
	message = ""
	mes := strings.ToLower(m.Content)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		//prevPost = m.MessageReference.MessageID
		//prevChan = m.ChannelID
		return
	}
	switch mes {
	case "!clear":
		Iclear()
	case "!list":
		message = ListGames()
	case "!trout":
		message = "trout that"
	case "!roll":
		message = IRoll()
	}

	if len(mes) > 5 {
		switch mes[:5] {
		case "!rmp ":
			message = Rmp(mes)
		case "!rmv ":
			message = Rmv(mes)
		case "!pick":
			message = IPick(mes, m)
		case "!veto":
			message = IVeto(mes, m)
		}
	}

	if message != "" {
		// If there's a previous message, delete it
		if previousMessageID != "" {
			err = s.ChannelMessageDelete(m.ChannelID, previousMessageID)
			if err != nil {
				fmt.Println("error deleting message,", err)
			}
		}
		// Send a text message
		msg, err := s.ChannelMessageSend(m.ChannelID, message)

		if err != nil {
			fmt.Println(err)
		} else {
			previousMessageID = msg.ID
		}

	} else {
		return
	}
}
