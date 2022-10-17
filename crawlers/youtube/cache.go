package youtube

import (
	"fmt"

	"github.com/eric2788/PlatformsCrawler/crawling"
)

type cache struct {
	prefix string
}

func (c *cache) key(channel string) string {
	return fmt.Sprintf("youtube:%s:%s", c.prefix, channel)
}

func (c *cache) setStruct(channel string, arg interface{}) error {
	return crawling.Store(c.key(channel), arg)
}

func (c *cache) getStruct(channel string, res interface{}) bool {
	success, err := crawling.GetStruct(c.key(channel), res)

	if !success && err != nil {
		logger.Errorf("從 redis 獲取 %s 時出現錯誤: %v", c.prefix, err)
	}

	return success
}

func (c *cache) getString(channel string) (string, bool) {
	result, err := crawling.GetString(c.key(channel))
	if err != nil {
		logger.Errorf("從 redis 獲取 %s 時出現錯誤: %v", c.prefix, err)
	}
	if result != "" {
		return result, true
	} else {
		return "", false
	}
}

func (c *cache) setString(channel, value string) error {
	return crawling.SetString(c.key(channel), value)
}

func (c *cache) GetAsLiveStatus(channel string) (*ChannelStatus, bool) {
	liveStatus := &ChannelStatus{}
	ok := c.getStruct(channel, liveStatus)
	return liveStatus, ok
}
