package youtube

import (
	"bufio"
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

	defer res.Body.Close()

	content, err := readAndFind(res, idRegex)

	if err != nil {
		return nil, err
	}

	if content == "" {
		return &ChannelStatus{Type: None}, nil // no streaming or upcoming
	} else {

		find := regex.GetParams(idRegex, content)
		videoId := find["id"]

		status := &ChannelStatus{Id: videoId}

		content, err = readAndFind(res, upcomingRegex)

		if err != nil {
			return nil, err
		}

		if content != "" {
			status.Type = UpComing
		} else {
			status.Type = Live
		}

		return status, nil
	}
}

// readAndFind 每行搜索，遇到符合條件立刻返回，防止搜索過長的字符串
func readAndFind(res *http.Response, reg *regexp.Regexp) (string, error) {
	bufReader := bufio.NewReader(res.Body)
	for line, _, err := bufReader.ReadLine(); err != io.EOF; line, _, err = bufReader.ReadLine() {
		if err != nil {
			return "", err
		}
		if reg.MatchString(string(line)) {
			return string(line), nil
		}
	}
	return "", nil
}
