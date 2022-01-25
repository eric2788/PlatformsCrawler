package youtube

import (
	"fmt"
	"github.com/eric2788/common-utils/regex"
	"io"
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

	content, err := readAsString(res)

	if err != nil {
		return nil, err
	}

	if !idRegex.MatchString(content) {
		return &ChannelStatus{Type: None}, nil // no streaming or upcoming
	} else {

		find := regex.GetParams(idRegex, content)

		videoId := find["id"]

		status := &ChannelStatus{Id: videoId}

		if upcomingRegex.MatchString(content) {
			status.Type = UpComing
		} else {
			status.Type = Live
		}

		return status, nil
	}
}

func readAsString(res *http.Response) (string, error) {

	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Warnf("response body close error: %v", err)
		}
	}()

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
