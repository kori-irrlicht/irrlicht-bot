package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var token string

func init() {
	flag.StringVar(&token, "t", "", "Discord authentication token")
}

func main() {
	flag.Parse()

	if token == "" {
		fmt.Println("No token")
		os.Exit(1)
	}
	Session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Can't create session", err)
	}

	Session.State.User, err = Session.User("@me")
	if err != nil {
		fmt.Println("Error fetching user info", err)
	}

	err = Session.Open()
	if err != nil {
		fmt.Println("Error connecting to discord", err)
	}

	//Session.ChannelMessageSend("274304138953752578", "Hallo Welt")
	channels, _ := Session.GuildChannels("208374693185585152")
	for _, channel := range channels {
		if channel.Type == "text" {
			Session.ChannelMessageSend(channel.ID, fmt.Sprintf("Dies ist der Channel %s.", channel.Name))
		}
	}

	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()

}
