package bot

import "net/http"
import "encoding/json"
import "io/ioutil"

const twitchAPIEndPoint = "https://api.twitch.tv/kraken"
const twitchChannelURLPrefix = "https://www.twitch.tv/"

type twitchStreams struct {
	Stream twitchStream `json:"stream"`
}

type twitchStream struct {
	ID      string        `json:"_id"`
	Game    string        `json:"game"`
	Viewers int           `json:"viewers"`
	Channel twitchChannel `json:"channel"`
}

type twitchChannel struct {
	Status string `json:"status"`
}

//isTwitchChannelOnline returns a Twitch channel's online status.
func isTwitchChannelOnline(channelName string) (isOnline bool, err error) {
	resp, err := http.Get(twitchAPIEndPoint + "/streams/" + channelName)
	if err != nil {
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

	if twitchStreams.Stream.ID != "" {
		return true, nil
	}

	return
}
