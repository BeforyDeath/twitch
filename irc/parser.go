package irc

import (
	"fmt"
	"strings"
)

type message struct {
	Channel string
	Origin  string
	Command string
	Options []string
	Params  map[string]string
	Text    string

	sections []string
	index    index
}

type index struct {
	message int
	command int
	params  int
}

func Parser(msg string) message {
	if msg != "" {
		msg = strings.Replace(msg, "\r\n", "", -1)
		sections := strings.SplitN(msg, " :", 3)

		// todo частный случай для NOTICE, USERNOTICE, ROOMSTATE, когда текст упущен
		// todo fixed PING
		if sections[0] != "PING" && len(sections) == 2 && strings.HasPrefix(sections[1], "tmi.twitch.tv") {
			sections = append(sections, "")
		}

		m := message{sections: sections}

		m.GetIndexes()
		m.GetCommand()
		m.GetParams()
		m.GetMessage()

		return m
	}
	return message{}
}

func (m *message) GetIndexes() {
	m.index.message, m.index.command, m.index.params = 1, 0, 1
	if len(m.sections) > 2 {
		m.index.message, m.index.command, m.index.params = 2, 1, 0
	}
}

func (m *message) GetCommand() {
	f := strings.Fields(m.sections[m.index.command])
	switch len(f) {
	case 1:
		m.Command = f[0]
	case 2:
		m.Origin = strings.FieldsFunc(f[0], func(c rune) bool { return !(c != '!' && c != ':') })[0]
		m.Command = f[1]
	default:
		m.Origin = strings.FieldsFunc(f[0], func(c rune) bool { return !(c != '!' && c != ':') })[0]
		m.Command = f[1]
		m.Channel = f[2]
		m.Options = f[3:]
	}

	// todo bag {"Channel":"","Origin":"","Command":"tmi.twitch.tv","Options":null,"Params":{},"Text":""}
	if m.Command == "tmi.twitch.tv" {
		fmt.Printf("%q\n-------------------------------------\n", m.sections)
	}
}

func (m *message) GetParams() {
	m.Params = make(map[string]string)

	if len(m.sections) > m.index.params {
		params := m.sections[m.index.params]

		if strings.HasPrefix(params, "@") {
			s := strings.Replace(params, "@", "", 1)
			p := strings.Split(s, ";")
			for _, d := range p {
				pv := strings.Split(d, "=")
				m.Params[pv[0]] = pv[1]
			}
		}
	}
}

func (m *message) GetMessage() {
	if len(m.sections) > m.index.message {
		m.Text = m.sections[m.index.message]
	}
}
