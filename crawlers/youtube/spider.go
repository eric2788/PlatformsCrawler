package youtube

import (
	"context"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"google.golang.org/api/youtube/v3"
	"sync"
	"time"
)

const channelNameKey = "youtube:channelNames"

var (
	statusMap      = &sync.Map{}
	lastLiveMap    = &sync.Map{}
	lastPendingMap = &sync.Map{}
)

type (
	LiveStatus string

	LiveBroadcast struct {
		ChannelId   string     `json:"channelId"`
		ChannelName string     `json:"channelName"`
		Status      LiveStatus `json:"status"`
		Info        *LiveInfo  `json:"info"`
	}

	LiveInfo struct {
		Cover       *string `json:"cover"`
		Title       string  `json:"title"`
		Id          string  `json:"id"`
		PublishTime string  `json:"publishTime"`
		Description string  `json:"description"`
	}
)

func (e EventType) ToLiveStatus() LiveStatus {
	if e == None || e == Completed {
		return "idle"
	}
	return LiveStatus(e)
}

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

	statusMap.Store(channelId, &ChannelStatus{Type: None}) // init first state
	ticker := time.NewTicker(time.Second * time.Duration(youtubeYaml.Interval))

	defer wg.Done()
	defer logger.Infof("頻道 %s 監控中止。", channelId)
	defer ticker.Stop()
	defer instance.channels.Remove(channelId)

	channelName := instance.getChannelName(channelId)

	logger.Infof("頻道 %s 監控已啟動。", channelName)

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

			logger.Debugf("%s 的狀態是 %v", channelName, status.Type)

			// 與上一次的狀態相同
			if lastStatus, ok := statusMap.Load(channelId); ok {
				last := lastStatus.(*ChannelStatus)
				if last.Id == status.Id && status.Type == last.Type {
					continue
				}
			}
			statusMap.Store(channelId, status)

			if status.Type == Live {
				// 上一次直播的 video id 跟本次相同
				if lastId, ok := lastLiveMap.Load(channelId); ok {
					if lastId.(string) == status.Id {
						continue
					}
				}
				lastLiveMap.Store(channelId, status.Id)
			}

			if status.Type == UpComing {
				// 上一次預告的 video id 跟本次相同
				if lastId, ok := lastPendingMap.Load(channelId); ok {
					if lastId.(string) == status.Id {
						continue
					}
				}
				lastPendingMap.Store(channelId, status.Id)
			}
			go handleBroadcast(channelId, status, publisher)
		}
	}
}

func handleBroadcast(channelId string, status *ChannelStatus, publisher crawling.Publisher) {

	name := instance.getChannelName(channelId)

	broadcast := &LiveBroadcast{
		Status:      status.Type.ToLiveStatus(),
		ChannelId:   channelId,
		ChannelName: name,
	}

	defer publisher(channelId, broadcast)

	switch status.Type {
	case UpComing:
		logger.Infof("%s 在油管有預定直播。", name)
		break
	case Live:
		logger.Infof("%s 正在油管直播。", name)
		break
	default:
		logger.Infof("%s 的油管直播已結束。", name)
		return
	}

	// only upcoming and live can get video
	videos, err := getVideos(status.Id)
	if err != nil {
		logger.Errorf("嘗試獲取油管視頻資訊 %s 時出現錯誤: %v", status.Id, err)
		return
	} else if len(videos) == 0 {
		logger.Warnf("找不到 %s 的油管視頻 %s", name, status.Id)
		return
	}

	video := videos[0]

	broadcast.Info = &LiveInfo{
		Cover:       getCover(video.Snippet.Thumbnails),
		Title:       video.Snippet.Title,
		Id:          video.Id,
		PublishTime: getPublishTime(video),
		Description: video.Snippet.Description,
	}
}

func getCover(details *youtube.ThumbnailDetails) *string {
	switch {
	case details.Maxres != nil:
		return &details.Maxres.Url
	case details.Standard != nil:
		return &details.Standard.Url
	case details.High != nil:
		return &details.High.Url
	case details.Medium != nil:
		return &details.Medium.Url
	case details.Default != nil:
		return &details.Default.Url
	default:
		return nil
	}
}

func getPublishTime(video *youtube.Video) string {

	var publishTime string
	if video.LiveStreamingDetails != nil {
		d := video.LiveStreamingDetails
		if d.ActualStartTime != "" {
			publishTime = d.ActualStartTime
		} else {
			publishTime = d.ScheduledStartTime
		}
	} else {
		publishTime = video.Snippet.PublishedAt
	}
	return publishTime
}
