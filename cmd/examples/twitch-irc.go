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
	rooms.Add("etozhemad", "c_a_k_e", "dreadztv", "mistafaker", "etozhezanuda", "mob5tertv", "guit88man")
	err = rooms.Add("tsm_doublelift", "dreamhackcs", "test example error room", "garenatw", "dreamhackoverwatch")
	if err != nil {
		fmt.Println(err)
	}

	Viewers := rooms.Get("c_a_k_e").Viewers
	fmt.Printf("%v\n", Viewers)

	rooms.Join(TwitchIRCNick)
	rooms.Join("etozhemad", "c_a_k_e", "dreadztv", "mistafaker", "etozhezanuda", "mob5tertv", "guit88man")
	rooms.Join("tsm_doublelift", "dreamhackcs", "garenatw", "dreamhackoverwatch")

	go func() {
		for range time.Tick(time.Second * 10) {
			rooms.Send(TwitchIRCNick, "Test send message")
			rooms.Leave(TwitchIRCNick)
		}
	}()

	go func() {
		for {
			msg := ch.ReadNext()

			m := irc.Parser(msg)

			switch m.Command {
			case "PING":
				ch.Pong()
			/*/
			case "001": // The connection is established
			case "353": // The list of current chatters in a channel /NAMES list
			case "366": // End of /NAMES list
			case "421": // Unknown command
			case "002", "003", "004", "372", "375", "376": // Ignore
			case "CAP": // The client capability negotiation extension
			/*/
			case "JOIN": // Someone joined a channel
				if m.Origin == TwitchIRCNick {
					rooms.Joined(m.Channel, true)
				}

				obj, _ := json.Marshal(m)
				fmt.Printf("> %s\n", obj)

			case "PART": // Someone left a channel
				if m.Origin == TwitchIRCNick {
					rooms.Joined(m.Channel, false)
				}

				obj, _ := json.Marshal(m)
				fmt.Printf("> %s\n", obj)
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
				fmt.Printf("> %v:%v > %v\n", m.Channel, m.Origin, m.Text)

			default:
				obj, _ := json.Marshal(m)
				fmt.Printf("> %s\n", obj)
			}
		}
	}()

	quit := make(chan bool)
	<-quit
}
