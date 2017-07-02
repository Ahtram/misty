package bot

import "net/http"
import "encoding/json"
import "io/ioutil"

const beamAPIEndPoint = "https://beam.pro/api/v1"
const beamChannelURLPrefix = "https://beam.pro/"

type beamChannels struct {
	Online bool
}

//[Deprecated]
//isBeamChannelOnline returns a Beam channel's online status.
func isBeamChannelOnline(channelName string) (isOnline bool, err error) {
	resp, err := http.Get(beamAPIEndPoint + "/channels/" + channelName)
	if err != nil {
		return
	}

	respByteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	beamChannels := beamChannels{}
	err = json.Unmarshal(respByteArray, &beamChannels)
	if err != nil {
		return
	}

	if beamChannels.Online {
		return true, nil
	}

	return
}
