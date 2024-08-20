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
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()

	logger := &Logger{
		file: logFile,
	}
	logger.Println(fmt.Sprintf("Bot initiated on ver %s", gitDescription))

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		logger.Println(fmt.Sprintf("error creating Discord session, %e", err))
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, logger)
	})

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		logger.Println(fmt.Sprintf("error opening connection, %e", err))
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	logger.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, logger *Logger) {
	// This function will be called every time a new message is created on any channel that the authenticated bot has access to.
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

	// Log the incoming message
	logger.Println(fmt.Sprintf("Received message: %s, from %s, on %s", m.Content, m.Author.Username, m.ChannelID))

	// Look for basic commands
	switch mes {
	case "!clear":
		Iclear()
	case "!list":
		message = ListGames()
	case "!roll":
		message = IRoll(m, logger)
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

	// Return early if no message was created by any commands
	if message == "" {
		return
	}

	// Log the response message
	logger.Println(fmt.Sprintf("Sending message: %s", message))

	// If there's a previous message id saved, and we aren't keeping the current message, delete
	if previousMessageID != "" && !keepMessage {
		err = s.ChannelMessageDelete(m.ChannelID, previousMessageID)
		if err != nil {
			logger.Println(fmt.Sprintf("Error deleting message: %v", err))
		}
	}

	// Send a text message
	msg, err := s.ChannelMessageSend(m.ChannelID, message)

	// Record the message ID to delete on the next post... or don't record it if this message should be kept
	if err != nil {
		logger.Println(fmt.Sprintf("Error sending message: %v", err))
	} else if keepMessage {
		previousMessageID = ""
	} else {
		previousMessageID = msg.ID
	}
}
