package irc

import (
	"errors"
	"fmt"
)

type rooms struct {
	ch    *channel
	items map[string]*room
}

type room struct {
	title    string
	Chatters int
	Viewers  int
	joined   bool
}

func NewRooms(ch *channel) rooms {
	r := rooms{
		ch:    ch,
		items: make(map[string]*room),
	}
	return r
}

func (r rooms) Add(titles ...string) error {
	var not_exist []string

	for _, t := range titles {
		fmt.Printf("Add room %v : ", t)
		s, err := r.ch.api.GetStream(t)
		if err != nil {
			not_exist = append(not_exist, t)
			fmt.Println("no")
		} else {
			fmt.Println("yes")
		}
		r.items[t] = &room{
			title:   t,
			Viewers: s.Stream.Viewers}
	}

	if len(not_exist) > 0 {
		return errors.New(fmt.Sprintf("The rooms does not exist: %q", not_exist))
	}
	return nil
}

func (r rooms) Delete(title string) {

}

func (r rooms) Join(titles ...string) {
	for _, t := range titles {
		if room, ok := r.items[t]; ok && !room.joined {
			fmt.Printf("< JOIN #%v\r\n", t)
			fmt.Fprintf(r.ch.conn, "JOIN #%v\r\n", t)
		}
	}
}

func (r rooms) Leave(title string) {
	if room, ok := r.items[title]; ok && room.joined {
		fmt.Printf("< PART #%v\r\n", title)
		fmt.Fprintf(r.ch.conn, "PART #%v\r\n", title)
	}
}

func (r rooms) Joined(ch string, b bool) {
	r.items[ch].joined = b
}

func (r rooms) Send(title, msg string) {
	if room, ok := r.items[title]; ok && room.joined {
		fmt.Printf("< PRIVMSG #%v :%v\r\n", title, msg)
		fmt.Fprintf(r.ch.conn, "PRIVMSG #%v :%v\r\n", title, msg)
	}
}

func (r rooms) Get(title string) *room {
	if room, ok := r.items[title]; ok {
		return room
	}
	return &room{}
}
