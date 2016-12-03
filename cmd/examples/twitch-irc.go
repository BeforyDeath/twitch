package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/BeforyDeath/twitch/irc"
)

const (
	// Twitch Chat OAuth Token Generator - http://twitchapps.com/tmi/
	TwitchIRCPassword string = ""
	TwitchIRCNick     string = ""
	TwitchAPIClientID string = ""
)

func main() {

	ch, err := irc.NewChannel(TwitchIRCPassword, TwitchIRCNick, TwitchAPIClientID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()

	rooms := irc.NewRooms(&ch)

	rooms.Add(TwitchIRCNick)

	go func() {
		for t := range time.Tick(time.Second * 30) {
			rooms.Send(TwitchIRCNick, fmt.Sprintf("tick %d:%d", t.Minute(), t.Second()))
		}
	}()

	rooms.Join(TwitchIRCNick)

	go func() {
		for {
			msg := ch.ReadNext()

			m := irc.Parser(msg)

			obj, _ := json.Marshal(m)
			fmt.Printf("> %s\n", obj)

			switch m.Command {
			case "PING":
				ch.Pong()

			/*/
			case "001": // The connection is established
			/*/
			case "353":
				// The list of current chatters in a channel /NAMES list
				// If there are greater than 1000 chatters in a room,
				// NAMES will only return the list of OPs currently in the room
			/*/
			case "366": // End of /NAMES list
			case "421": // Unknown command
			case "002", "003", "004", "372", "375", "376": // Ignore
			case "CAP": // The client capability negotiation extension
			/*/
			case "JOIN": // Someone joined a channel
				if m.Origin == TwitchIRCNick {
					rooms.Joined(m.Room, true)
				}

			case "PART": // Someone left a channel
				if m.Origin == TwitchIRCNick {
					rooms.Joined(m.Room, false)
				}

			/*/
			case "MODE": // Someone gained or lost operator
			case "CLEARCHAT": // Username is timed out or banned on a channel
			case "USERSTATE": // Is sent when joining a channel and every time you send a PRIVMSG to a channel
			case "ROOMSTATE": // is sent when joining a channel and every time one of the chat room settings, like slow mode, change
			case "GLOBALUSERSTATE": // is sent on successful login, if the capabilities have been acknowledged before then
			case "NOTICE": // General notices from the server - could be about state change
			case "USERNOTICE": // Re-subscription notice
			case "HOSTTARGET": // Host starts message
			case "RECONNECT": // Twitch IRC processes occasionally need to be restarted
			//*/
			case "PRIVMSG": // Message
				fmt.Printf("> %v:%v\t: %v\n", m.Room, m.Origin, m.Text)
			default:
				//fmt.Println("No command")
			}
		}
	}()

	quit := make(chan bool)
	<-quit
}
