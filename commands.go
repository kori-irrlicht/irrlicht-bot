package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const commandPrefix = `$`

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	user, _ := s.User("@me")
	if m.Author.ID != user.ID {
		s.ChannelMessageSend(botTestChannel, fmt.Sprintf("%s schrieb: %s", m.Author.Username, m.Content))

		if !strings.HasPrefix(m.Content, commandPrefix) {
			return
		}

		strAr := strings.Split(m.Content, " ")
		str := strings.Replace(strAr[0], commandPrefix, "", 1)

		switch str {
		case "help":
			s.ChannelMessageSend(m.ChannelID, "Hallo")
		}

	}
}
