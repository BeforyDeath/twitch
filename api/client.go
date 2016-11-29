package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

const (
	TwitchAPIHost = "https://api.twitch.tv/kraken/"
)

type Client struct {
	clientID string
	client   http.Client
}

func NewClient(clientID string) Client {
	return Client{clientID: clientID}
}

func (c Client) get(method string, result interface{}) error {
	req, err := http.NewRequest("GET", TwitchAPIHost+method, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.twitchtv.v3+json")
	if c.clientID != "" {
		req.Header.Add("Client-ID", c.clientID)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return errors.New("Twitch API status code: " + strconv.Itoa(resp.StatusCode))
	}

	// todo response debugging
	//tmp, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%v\n", string(tmp))
	//resp.Body = ioutil.NopCloser(bytes.NewBuffer(tmp))

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) GetStream(title string) (streams, error) {
	stream := new(streams)

	err := c.get("streams/"+title, stream)
	if err != nil {
		return streams{}, err
	}
	return *stream, nil
}
