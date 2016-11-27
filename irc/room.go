package irc

import (
	"fmt"
	"net"
)

// todo запросить зрителей (viewers int) и онлайн (onlain bool)
type room struct {
	conn    *net.Conn
	title   string
	onlain  bool
	viewers int
	joined  bool
}

// todo ждать подтвержение, что мы вошли в комнату (joined bool)
func (r room) Join() {
	if r.conn == nil {
		fmt.Println("Join nil connect")
		return
	}
	fmt.Fprintf(*r.conn, "JOIN #%v\r\n", r.title)
}

// todo ждать подтвержение, что мы выйшли из комнату (joined bool)
func (r room) Leave() {
	if r.conn == nil {
		fmt.Println("Leave nil connect")
		return
	}
	fmt.Fprintf(*r.conn, "PART #%v\r\n", r.title)
}

// todo проверить, действительно ли мы в комнате (joined bool)
func (r room) Send(msg string) {
	if r.conn == nil {
		fmt.Println("Send nil connect")
		return
	}

	if r.joined {
		fmt.Fprintf(*r.conn, "PRIVMSG #%v :%v\r\n", r.title, msg)
	}
}
