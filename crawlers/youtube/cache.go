package youtube

import (
	"github.com/eric2788/PlatformsCrawler/crawling"
)


func newCache(prefix string) *youtubeCache {
	return forYoutube(crawling.NewCache("youtube", prefix))
}

func forYoutube(c *crawling.Cache) *youtubeCache {
	return &youtubeCache{c}
}


type youtubeCache struct {
	*crawling.Cache
}

func (c *youtubeCache) GetAsLiveStatus(channel string) (*ChannelStatus, bool) {
	liveStatus := &ChannelStatus{}
	ok, err := c.GetStruct(channel, liveStatus)
	if err != nil {
		logger.Errorf("從 redis 獲取 %s 時出現錯誤: %v, 將嘗試使用本地快取", c.Prefix, err)
		if status, ok := c.Local[channel]; ok {
			return status.(*ChannelStatus), ok
		}
	}
	return liveStatus, ok
}
