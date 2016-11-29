package irc

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/BeforyDeath/twitch/api"
)

const (
	TwitchIRCHost string = "irc.twitch.tv:6667"
)

type channel struct {
	api    api.Client
	conn   net.Conn
	reader *bufio.Reader
}

func NewChannel(pass, nick, clientID string) (channel, error) {
	conn, err := net.Dial("tcp", TwitchIRCHost)
	if err != nil {
		return channel{}, err
	}

	fmt.Fprintf(conn, "PASS %s\r\nNICK %s\r\n", pass, nick)

	ch := channel{
		api:    api.NewClient(clientID),
		conn:   conn,
		reader: bufio.NewReader(conn),
	}

	ch.capabilities()
	return ch, nil
}

func (ch channel) Close() {
	ch.conn.Close()
}

func (ch channel) capabilities() {
	fmt.Fprint(ch.conn, "CAP REQ :twitch.tv/tags\r\n")
	time.Sleep(time.Millisecond * 100)
	fmt.Fprint(ch.conn, "CAP REQ :twitch.tv/commands\r\n")
	time.Sleep(time.Millisecond * 100)
	fmt.Fprint(ch.conn, "CAP REQ :twitch.tv/membership\r\n")
	time.Sleep(time.Millisecond * 100)
}

func (ch channel) ReadNext() string {
	msg, err := ch.reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return msg
}

func (ch channel) Pong() {
	fmt.Println("< PONG :tmi.twitch.tv")
	fmt.Fprint(ch.conn, "PONG :tmi.twitch.tv\r\n")
}
