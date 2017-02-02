package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"github.com/yasvisu/gw2api"
)

var botTestChannel string
var discordToken string
var guildWars2Token string

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	session, err := createDiscordSession()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	api, err := createGuildWarsSession()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	_ = api

	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	session.Close()

}

func createDiscordSession() (*discordgo.Session, error) {
	discordToken = viper.GetString("discordToken")
	if discordToken == "" {
		return nil, fmt.Errorf("No discordToken")
	}
	Session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Can't create session", err)
		return nil, err
	}

	Session.State.User, err = Session.User("@me")
	if err != nil {
		fmt.Println("Error fetching user info", err)
		return nil, err
	}

	err = Session.Open()
	if err != nil {
		fmt.Println("Error connecting to discord", err)
		return nil, err
	}

	Session.AddHandler(messageCreateHandler)
	return Session, nil
}

func createGuildWarsSession() (*gw2api.GW2Api, error) {
	guildWars2Token = viper.GetString("guildWars2Token")
	api, err := gw2api.NewAuthenticatedGW2Api(guildWars2Token)
	if err != nil {
		fmt.Println("Can't create gw2api")
		return nil, err
	}

	return api, nil
}
