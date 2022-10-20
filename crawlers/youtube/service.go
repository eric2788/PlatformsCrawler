package youtube

import (
	"context"
	"fmt"
	"os"

	"github.com/eric2788/PlatformsCrawler/crawling"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type EventType string

const (
	UpComing  EventType = "upcoming"
	Live      EventType = "live"
	Completed EventType = "completed"
	None      EventType = "none"
)

var (
	ctx            = context.Background()
	youtubeService *youtube.Service
)

func initYoutubeService() {
	service, err := youtube.NewService(ctx, option.WithAPIKey(youtubeYaml.Api.Key))
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
	youtubeService = service
}

// searchVideo -100 quota per request!
func searchVideo(channelId string, eventType EventType) (*youtube.SearchResultSnippet, error) {
	call := youtubeService.Search.
		List([]string{"snippet"}).
		ChannelId(channelId).
		EventType(string(eventType)).
		Type("video").
		Order("date").
		MaxResults(1).
		RegionCode(youtubeYaml.Api.Region).
		RelevanceLanguage(youtubeYaml.Api.Language)
	res, err := call.Do()
	if err != nil {
		return nil, err
	} else if len(res.Items) == 0 {
		return nil, nil
	} else {
		return res.Items[0].Snippet, nil
	}
}

// getChannels -1 quota per request
// return only snippet
func getChannels(channelId ...string) ([]*youtube.Channel, error) {
	call := youtubeService.Channels.
		List([]string{"snippet"}).
		Id(channelId...)
	channels := make([]*youtube.Channel, 0)
	err := call.Pages(ctx, func(res *youtube.ChannelListResponse) error {
		for _, item := range res.Items {
			channels = append(channels, item)
		}
		return nil
	})
	return channels, err
}

// getChannelByUsername -1 quota per quest but only can request once
// return only snippet
func getChannelByUsername(username string) (*youtube.ChannelSnippet, error) {
	call := youtubeService.Channels.
		List([]string{"snippet"}).
		ForUsername(username).
		MaxResults(1)
	res, err := call.Do()
	if err != nil {
		return nil, err
	} else if len(res.Items) == 0 {
		return nil, nil
	} else {
		return res.Items[0].Snippet, nil
	}
}

// getVideos -1 quota
// return snippet and liveStreamingDetails (if video is an archive)
func getVideos(id ...string) ([]*youtube.Video, error) {
	call := youtubeService.Videos.
		List([]string{"snippet", "liveStreamingDetails"}).
		Id(id...).
		RegionCode(youtubeYaml.Api.Region)

	videos := make([]*youtube.Video, 0)
	err := call.Pages(ctx, func(res *youtube.VideoListResponse) error {
		for _, item := range res.Items {
			videos = append(videos, item)
		}
		return nil
	})
	return videos, err
}

// getVideoWithCache with low consome quota (with redis cache)
func getVideoWithCache(id string, eType EventType) (*youtube.Video, error) {
	if eType != Live && eType != UpComing {
		return nil, fmt.Errorf("unsupported event type on getting video %s: %v", id, eType)
	}

	key := fmt.Sprintf("youtube:video_info_%s:%s", string(eType), id)
	video := &youtube.Video{}

	if ok, err := crawling.GetStruct(key, video); err != nil || !ok {

		if err != nil {
			logger.Errorf("從redis獲取在狀態 %v 的視頻ID %v 的資訊失敗: %v, 將使用 youtube data api.", eType, id, err)
		} else {
			logger.Warnf("無法在redis找到狀態 %v 的視頻ID %v 的資訊，將使用 youtube data api", eType, id)
		}

		fetch, err := getVideos(id)
		if err != nil {
			return nil, err
		}

		video = fetch[0]

		// save to redis
		if err := crawling.Store(key, video); err != nil {
			logger.Warnf("嘗試儲存狀態 %v 的視頻ID %v 的資訊到redis失敗: %v", eType, id, err)
		} else {
			logger.Infof("儲存狀態 %v 的視頻ID %v 的資訊到redis成功。", eType, id)
		}

	}

	return video, nil
}
