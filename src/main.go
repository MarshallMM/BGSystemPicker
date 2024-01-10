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
	lockdownEnabled bool
	authorizedList []string
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
		authorized bool
	)
	message = ""
	mes := strings.ToLower(m.Content)


	// Iterate through the slice to check if any value equals the specific one
	authorized = false
	for _, str := range authorizedList {
		if str == m.Author.ID {
			authorized = true
			break
		}
	}
	
	// Ignore all messages created by the bot itself or non auth users if lockdown is in place
	if m.Author.ID == s.State.User.ID {
		//prevPost = m.MessageReference.MessageID
		//prevChan = m.ChannelID
		return
	} else if !authorized && lockdownEnabled {
		return
	}

	// Authorized users can lockdown bot to only authorized users
	if authorized && mes == "!lockdown" {
		lockdownEnabled = !lockdownEnabled
		return
	}
		
	
	if authorized || !lockdownEnabled {	
		switch mes {
		case "!clear":
			Iclear()
		case "!list":
			message = ListGames()
		case "!trout":
			message = "trout that"
		case "!roll":
			message, err = IRoll()
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
	} else {
		fmt.Printf("%s attempted to use the below command: %s,m.Author.ID, m.content)
	}

	if message != "" {
		//err = s.ChannelMessageDelete(prevChan, prevPost)
		if err != nil {
			fmt.Println(err)
		}
		// Send a text message
		_, err = s.ChannelMessageSend(m.ChannelID, message)

		if err != nil {
			fmt.Println(err)
		}

	} 
}
