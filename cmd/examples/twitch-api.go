package main

import (
	"fmt"

	"github.com/BeforyDeath/twitch/api"
)

// Get Client-ID - register your application from the https://www.twitch.tv/settings/connections
const TwitchAPIClientID = ""

func main() {

	client := api.NewClient(TwitchAPIClientID)

	stream, err := client.GetStream("c_a_k_e")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", stream)
}
