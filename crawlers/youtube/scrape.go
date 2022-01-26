package youtube

import (
	"fmt"
	"github.com/eric2788/common-utils/regex"
	"github.com/eric2788/common-utils/request"
	"net/http"
	"regexp"
)

var (
	idRegex       *regexp.Regexp
	upcomingRegex *regexp.Regexp
)

func initKeywordRegexp() {
	idRegex = regexp.MustCompile(youtubeYaml.LiveKeyword)
	upcomingRegex = regexp.MustCompile(youtubeYaml.UpComingKeyword)
}

type ChannelStatus struct {
	Type EventType
	Id   string
}

func GetChannelStatus(channelId string) (*ChannelStatus, error) {

	res, err := http.Get(fmt.Sprintf("https://youtube.com/channel/%s/live", channelId))

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("not found channel %s", channelId)
	} else if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	contents, err := request.ReadForRegex(res, idRegex, upcomingRegex)

	if err != nil {
		return nil, err
	}

	if contents[0] == "" {
		return &ChannelStatus{Type: None}, nil // no streaming or upcoming
	} else {

		find := regex.GetParams(idRegex, contents[0])
		videoId := find["id"]

		status := &ChannelStatus{Id: videoId}

		if contents[1] != "" {
			status.Type = UpComing
		} else {
			status.Type = Live
		}

		return status, nil
	}
}
