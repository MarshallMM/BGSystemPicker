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
	gitDescription    string
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
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
	keepMessage := false

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Look for basic commands
	switch mes {
	case "!clear":
		Iclear()
	case "!list":
		message = ListGames()
	case "!roll":
		message = IRoll()
		keepMessage = true
	case "!version":
		message = gitDescription + "\nhttps://github.com/MarshallMM/BGSystemPicker"
		keepMessage = true
	}

	// Look for more complex commands with payloads
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
	// Return Early if no message was created by any commands
	if message == "" {
		return
	}

	// If there's a previous message id saved, and we arent keeping the current message, delete
	if previousMessageID != "" && !keepMessage {
		err = s.ChannelMessageDelete(m.ChannelID, previousMessageID)
		if err != nil {
			fmt.Println("error deleting message,", err)
		}
	}
	// Send a text message
	msg, err := s.ChannelMessageSend(m.ChannelID, message)

	// Record the message ID to delete on the next post... or dont record it if this message should be kept
	if err != nil {
		fmt.Println(err)
	} else if keepMessage {
		previousMessageID = ""
	} else {
		previousMessageID = msg.ID
	}

}
