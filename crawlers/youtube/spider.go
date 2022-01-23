package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"sync"
	"time"
)

const channelNameKey = "youtube:channelNames"

var (
	statusMap = &sync.Map{}
)

func getChannelNames(channelId ...string) (map[string]string, error) {
	isNotFound := func(error) bool { return false }
	return crawling.LoadWithCache(channelNameKey, lookupNamesByChannelIds, isNotFound, channelId...)
}

func lookupNamesByChannelIds(channelIds []string) (map[string]string, error) {
	res, err := getChannels(channelIds...)
	if err != nil {
		return nil, err
	}
	names := make(map[string]string)
	for _, channel := range res {
		names[channel.Id] = channel.Snippet.Title
	}
	return names, nil
}

func runYoutubeSpider(ctx context.Context, channelId string, wg *sync.WaitGroup, publisher crawling.Publisher) {

	statusMap.Store(channelId, None) // init first state
	ticker := time.NewTicker(time.Second * time.Duration(youtubeYaml.Interval))

	defer wg.Done()
	defer logger.Infof("頻道 %s 監控中止。", channelId)
	defer ticker.Stop()
	defer instance.channels.Remove(channelId)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			status, err := GetChannelStatus(channelId)
			if err != nil {
				logger.Errorf("嘗試獲取頻道 %s 的直播狀態時出現錯誤: %v", channelId, err)
				continue
			}

			logger.Debugf("%s 的狀態是 %v", instance.getChannelName(channelId), status.Type)

			// 與上一次的狀態相同
			if lastStatus, ok := statusMap.Load(channelId); ok && lastStatus.(EventType) == status.Type {
				continue
			}

			statusMap.Store(channelId, status.Type)
			go handleBroadcast(channelId, status, publisher)
		}
	}
}

func handleBroadcast(channelId string, status *ChannelStatus, publisher crawling.Publisher) {

	name := instance.getChannelName(channelId)

	switch status.Type {
	case UpComing:
		logger.Infof("%s 在油管有預定直播: ", name)
		break
	case Live:
		logger.Infof("%s 正在油管直播: ", name)
		break
	default:
		logger.Infof("%s 的油管直播已結束。", name)
		return
	}

	// only upcoming and live can get video
	video, err := getVideos(status.Id)
	if err != nil {
		logger.Errorf("嘗試獲取油管視頻資訊 %s 時出現錯誤: %v", status.Id, err)
		return
	} else if len(video) == 0 {
		logger.Warnf("找不到 %s 的油管視頻 %s", name, status.Id)
		return
	}

	if b, err := json.MarshalIndent(video[0].Snippet, "", "\t"); err != nil {
		logger.Errorf("json error: %v", err)
	} else {
		fmt.Println(string(b))
	}
}
