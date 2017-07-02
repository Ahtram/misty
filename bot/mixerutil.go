package bot

import "net/http"
import "encoding/json"
import "io/ioutil"

const mixerAPIEndPoint = "https://mixer.com/api/v1/"
const mixerChannelURLPrefix = "https://mixer.com/"

type mixerChannels struct {
	Online bool
}

//isMixerChannelOnline returns a Beam channel's online status.
func isMixerChannelOnline(channelName string) (isOnline bool, err error) {
	resp, err := http.Get(mixerAPIEndPoint + "channels/" + channelName)
	if err != nil {
		return
	}

	respByteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	mixerChannels := mixerChannels{}
	err = json.Unmarshal(respByteArray, &mixerChannels)
	if err != nil {
		return
	}

	if mixerChannels.Online {
		return true, nil
	}

	return
}
