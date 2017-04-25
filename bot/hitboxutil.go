package bot

import "net/http"
import "encoding/json"
import "io/ioutil"

const hitboxAPIEndPoint = "https://api.hitbox.tv"
const hitboxChannelURLPrefix = "http://www.hitbox.tv/"

type hitboxChannels struct {
	Islive string `json:"Is_live"`
}

//isHitboxChannelOnline returns a hitbox channel's online status.
func isHitboxChannelOnline(channelName string) (isOnline bool, err error) {
	resp, err := http.Get(hitboxAPIEndPoint + "/user/" + channelName)
	if err != nil {
		return
	}

	respByteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	hitboxChannels := hitboxChannels{}
	err = json.Unmarshal(respByteArray, &hitboxChannels)
	if err != nil {
		return
	}

	if hitboxChannels.Islive == "1" {
		return true, nil
	}

	return
}
