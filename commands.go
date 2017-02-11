package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

const commandPrefix = `$`
const nameSplit = `#`
const regelChan = `regeln`

var globalCommands = map[string]Command{
	"help": &Help{},
	"name": &Name{},
}
var ruleCommands = map[string]Command{
	"accept": &Accept{},
	"update": &UpdateRules{},
}

/**
Command represents the available chat commands
*/
type Command interface {
	Exec(s *discordgo.Session, m *discordgo.MessageCreate)
	Description() string
}

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	user, _ := s.User("@me")
	if m.Author.ID != user.ID {
		ch, _ := s.Channel(m.ChannelID)
		if ch.Name == regelChan {
			s.ChannelMessageDelete(ch.ID, m.ID)
			return
		}
		if !strings.HasPrefix(m.Content, commandPrefix) {
			return
		}

		strAr := strings.Split(m.Content, " ")
		str := strings.Replace(strAr[0], commandPrefix, "", 1)

		comm, ok := globalCommands[str]
		if !ok {
			if ch.Name != regelChan {
				return
			}
			comm, ok = ruleCommands[str]
			if !ok {
				return
			}
		}
		comm.Exec(s, m)

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
		max := 0
		for i := 0; i < len(file); i += max {
			max = i + 2000
			if len(file) <= max {
				max = len(file)
			} else {
				max = strings.LastIndex(string(file[:max]), "\n")
			}
			_, err := s.ChannelMessageSend(rules, string(file[i:max]))
			if err != nil {
				return err
			}
		}
	}

	return nil

}

/**
Accept is the command to accept the rules
*/
type Accept struct {
	Command
}

func (a *Accept) Description() string {
	return `Akzeptiert die Regeln und weist dem User die 'Mitglied'-Rolle zu.`
}

func (a *Accept) Exec(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, _ := s.Channel(m.ChannelID)
	if ch.Name != regelChan {
		return
	}
	err := s.GuildMemberRoleAdd(ch.GuildID, m.Author.ID, viper.GetString("mitgliedRollenId"))
	if err != nil {
		fmt.Println(err)
	}
	s.ChannelMessageDelete(ch.ID, m.ID)

}

/**
UpdateRules removes all messages from the channel
and updates the rules */
type UpdateRules struct {
	Command
}

func (u *UpdateRules) Description() string {
	return `Aktualisiert die Regeln.`
}

func (u *UpdateRules) Exec(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, _ := s.Channel(m.ChannelID)
	if ch.Name != regelChan {
		return
	}
	msgs, err := s.ChannelMessages(ch.ID, 100, "999999999999999999", "000000000000000000")
	if err != nil {
		fmt.Println(err)
	}
	for _, msg := range msgs {
		s.ChannelMessageDelete(ch.ID, msg.ID)
	}
	checkRuleEntry(s)
}

type Help struct {
	Command
}

func (h *Help) Description() string {
	return `Gibt die verfügbaren Kommandos aus.`
}

func (h *Help) Exec(s *discordgo.Session, m *discordgo.MessageCreate) {
	printer := func(title string, cmdMap map[string]Command) string {
		m := fmt.Sprintf("**%s**\n", title)
		for k, v := range cmdMap {
			m += fmt.Sprintf(" - `$%s`: %s\n", k, v.Description())
		}
		return m
	}

	msg := printer("Globale Kommandos", globalCommands)
	msg += "\n"
	msg += printer("Kommandos für den Regelchannel", ruleCommands)
	s.ChannelMessageSend(m.ChannelID, msg)
}

type Name struct {
	Command
}

func (n *Name) Description() string {
	return `Registriere deinen GW2 Namen.`
}

func (n *Name) Exec(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, _ := s.Channel(m.ChannelID)
	fmt.Println(m.Content)
	if !strings.Contains(m.Content, " ") {
		s.ChannelMessageSend(m.ChannelID, "Kein Name angegeben.")
		return
	}
	newname := strings.SplitN(string(m.Content), " ", 2)[1]

	err := s.GuildMemberNickname(ch.GuildID, m.Author.ID, m.Author.Username+nameSplit+newname)
	if err != nil {
		fmt.Println(err)
	}
}
