package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

const commandPrefix = `$`

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	user, _ := s.User("@me")
	if m.Author.ID != user.ID {
		ch, _ := s.Channel(m.ChannelID)
		if ch.Name == "regeln" && !strings.HasPrefix(m.Content, commandPrefix) {
			s.ChannelMessageDelete(ch.ID, m.ID)
			return
		}
		if !strings.HasPrefix(m.Content, commandPrefix) {
			return
		}

		strAr := strings.Split(m.Content, " ")
		str := strings.Replace(strAr[0], commandPrefix, "", 1)
		if ch.Name == "regeln" {

			switch str {
			case "accept":
				rollen, _ := s.GuildRoles(ch.GuildID)
				for _, rl := range rollen {
					fmt.Println(rl.ID)
				}
				err := s.GuildMemberRoleAdd(ch.GuildID, m.Author.ID, viper.GetString("mitgliedRollenId"))
				if err != nil {
					fmt.Println(err)
				}
				s.ChannelMessageDelete(ch.ID, m.ID)

			case "update":
				// Remove all message from the channel
				// and update the rules
				msgs, err := s.ChannelMessages(ch.ID, 100, "999999999999999999", "000000000000000000")
				if err != nil {
					fmt.Println(err)
				}
				for _, msg := range msgs {
					s.ChannelMessageDelete(ch.ID, msg.ID)
				}
				checkRuleEntry(s)
			default:
				s.ChannelMessageDelete(ch.ID, m.ID)
			}

		}

		switch str {
		case "help":
			s.ChannelMessageSend(m.ChannelID, "Hallo")
		}

	}
}

func checkRuleEntry(s *discordgo.Session) (err error) {
	rules := viper.GetString("channel.rules")

	msgs, err := s.ChannelMessagesPinned(rules)
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(viper.GetString("regeln"))
	if err != nil {
		return err
	}

	if len(msgs) == 0 {
		for i := 0; i < len(file); i += 2000 {

			max := i + 2000
			if max > len(file) {
				max = len(file)
			}

			_, err := s.ChannelMessageSend(rules, string(file[i:max]))
			if err != nil {
				return err
			}
		}
	}

	return nil

}
