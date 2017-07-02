package bot

import "net/http"
import "encoding/json"
import "io/ioutil"

const smashcastAPIEndPoint = "https://api.smashcast.tv/"
const smashcastChannelURLPrefix = "https://www.smashcast.tv/"

type smashcastChannels struct {
	Islive string `json:"Is_live"`
}

//isSmashcastChannelOnline returns a smashcast channel's online status.
func isSmashcastChannelOnline(channelName string) (isOnline bool, err error) {
	resp, err := http.Get(smashcastAPIEndPoint + "user/" + channelName)
	if err != nil {
		return
	}

	respByteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	smashcastChannels := smashcastChannels{}
	err = json.Unmarshal(respByteArray, &smashcastChannels)
	if err != nil {
		return
	}

	if smashcastChannels.Islive == "1" {
		return true, nil
	}

	return
}
