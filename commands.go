package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

const commandPrefix = `$`

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	user, _ := s.User("@me")
	if m.Author.ID != user.ID {
		if !strings.HasPrefix(m.Content, commandPrefix) {
			return
		}

		strAr := strings.Split(m.Content, " ")
		str := strings.Replace(strAr[0], commandPrefix, "", 1)
		ch, _ := s.Channel(m.ChannelID)
		if ch.Name == "regeln" {

			switch str {
			case "accept":
				rollen, _ := s.GuildRoles(ch.GuildID)
				for _, rl := range rollen {
					fmt.Println(rl.ID)
				}
				fmt.Println(s.GuildRoles(ch.GuildID))
				err := s.GuildMemberRoleAdd(ch.GuildID, m.Author.ID, viper.GetString("mitgliedRollenId"))
				fmt.Println(err)
			}
		}

		switch str {
		case "help":
			s.ChannelMessageSend(m.ChannelID, "Hallo")
		}

	}
}
