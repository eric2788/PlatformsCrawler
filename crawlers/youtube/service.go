package youtube

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
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
