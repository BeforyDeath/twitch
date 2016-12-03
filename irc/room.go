package irc

import "fmt"

type rooms struct {
	ch    *channel
	items map[string]*room
}

type room struct {
	title  string
	online bool
	joined bool
}

func NewRooms(ch *channel) rooms {
	r := rooms{
		ch:    ch,
		items: make(map[string]*room),
	}
	return r
}

func (r rooms) Add(titles ...string) {
	for _, t := range titles {
		fmt.Printf("Add room %v \n", t)
		r.items[t] = &room{title: t}
	}
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

func (r rooms) GetAll() map[string]*room {
	return r.items
}
