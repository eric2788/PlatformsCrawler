package youtube

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/corpix/uarand"
	"github.com/eric2788/common-utils/regex"
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

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://youtube.com/channel/%s/live", channelId), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", uarand.GetRandom())

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	} else if res.StatusCode == 404 {
		return nil, fmt.Errorf("not found channel %s", channelId)
	}

	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	isUpcoming, hasConical := upcomingRegex.Match(content), idRegex.Match(content)

	if !hasConical {
		return &ChannelStatus{Type: None}, nil // no streaming or upcoming
	} else {

		find := regex.GetParams(idRegex, string(content))
		videoId := find["id"]

		status := &ChannelStatus{Id: videoId}

		if isUpcoming {
			status.Type = UpComing
		} else {
			status.Type = Live
		}

		return status, nil
	}
}
