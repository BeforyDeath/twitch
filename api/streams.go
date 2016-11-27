package api

import "time"

// ref https://dev.twitch.tv/docs/api/v3/streams/

type streams struct {
	Stream stream `json:"stream,omitempty"`
}

type stream struct {
	Game       string    `json:"game,omitempty"`
	Viewers    int       `json:"viewers,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"` // RUS +3:00
	Channel    channel   `json:"channel,omitempty"`
}

type channel struct {
	Status       string    `json:"status,omitempty"`
	Name         string    `json:"name,omitempty"`
	Display_name string    `json:"display_name,omitempty"`
	Updated_at   time.Time `json:"updated_at,omitempty"` // RUS +3:00
}
