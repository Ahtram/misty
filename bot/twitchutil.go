package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const twitchAPIEndPoint = "https://api.twitch.tv/kraken"
const twitchChannelURLPrefix = "https://www.twitch.tv/"
const twitchClientID = "9mfyew1nli07zwkxa6rcn5utcgplk2"

type twitchStreams struct {
	Stream twitchStream `json:"stream"`
}

type twitchStream struct {
	ID      int           `json:"_id"`
	Game    string        `json:"game"`
	Viewers int           `json:"viewers"`
	Preview twitchPreview `json:"preview"`
	Channel twitchChannel `json:"channel"`
}

type twitchChannel struct {
	Status string `json:"status"`
}

type twitchPreview struct {
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Large    string `json:"large"`
	Template string `json:"template"`
}

//isTwitchChannelOnline returns a Twitch channel's online status.
func isTwitchChannelOnline(channelName string) (isOnline bool, err error) {
	request, err := http.NewRequest("GET", twitchAPIEndPoint+"/streams/"+channelName, nil)
	if err != nil {
		fmt.Println(Red("[POST Request Error] ") + err.Error())
		return
	}

	request.Header.Set("Client-ID", twitchClientID)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(Red("[POST Request Error] ") + err.Error())
		return
	}

	respByteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	twitchStreams := twitchStreams{}
	err = json.Unmarshal(respByteArray, &twitchStreams)
	if err != nil {
		return
	}

	if twitchStreams.Stream.Channel.Status != "" {
		return true, nil
	}

	return
}
