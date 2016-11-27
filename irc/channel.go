package irc

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	TwitchIRCHost string = "irc.twitch.tv:6667"
)

type channel struct {
	conn   net.Conn
	reader *bufio.Reader
	rooms  map[string]room
}

func Connect(pass, nick string) (channel, error) {
	conn, err := net.Dial("tcp", TwitchIRCHost)
	if err != nil {
		return channel{}, err
	}

	fmt.Fprintf(conn, "PASS %s\r\nNICK %s\r\n", pass, nick)

	ch := channel{
		conn:   conn,
		reader: bufio.NewReader(conn),
		rooms:  make(map[string]room),
	}

	time.Sleep(time.Millisecond * 100)
	fmt.Fprint(conn, "CAP REQ :twitch.tv/tags\r\n")
	time.Sleep(time.Millisecond * 100)
	fmt.Fprint(conn, "CAP REQ :twitch.tv/commands\r\n")
	time.Sleep(time.Millisecond * 100)
	fmt.Fprint(conn, "CAP REQ :twitch.tv/membership\r\n")
	time.Sleep(time.Millisecond * 100)

	return ch, nil
}

func (ch *channel) Close() {
	ch.conn.Close()
}

func (ch *channel) Pong() {
	fmt.Fprint(ch.conn, "PONG :tmi.twitch.tv\r\n")
}

func (ch channel) ReadNext() string {
	msg, err := ch.reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return msg
}

// todo сделать запрос к api, существует ли комната
func (ch *channel) AddRooms(title ...string) int {
	for _, t := range title {
		ch.rooms[t] = room{
			conn:  &ch.conn,
			title: t,
		}
		ch.rooms[t].Join()
	}
	return len(title)
}

func (ch *channel) DeleteRooms(title ...string) int {
	for _, t := range title {
		ch.rooms[t].Leave()
		delete(ch.rooms, t)
	}
	return len(title)
}

func (ch *channel) Room(title string) room {
	if r, ok := ch.rooms[title]; ok {
		return r
	}
	return room{}
}

func (ch *channel) GetRoom(title string) (room, error) {
	if r, ok := ch.rooms[title]; ok {
		return r, nil
	}
	return room{}, errors.New("Error: No room #" + title)
}
