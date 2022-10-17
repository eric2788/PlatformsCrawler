package youtube

import (
	"fmt"

	"github.com/eric2788/PlatformsCrawler/crawling"
)

type cache struct {
	prefix string
	local map[string]interface{}
}

func newCache(prefix string) *cache {
	return &cache{
		prefix: prefix,
		local: make(map[string]interface{}),
	}
}

func (c *cache) key(channel string) string {
	return fmt.Sprintf("youtube:%s:%s", c.prefix, channel)
}

func (c *cache) setStruct(channel string, arg interface{}) error {
	err := crawling.Store(c.key(channel), arg)
	if err != nil {
		logger.Errorf("儲存 %s 到 redis 時出現錯誤: %v, 將嘗試使用本地快取", c.prefix, err)
		c.local[channel] = arg
	}
	return err
}

func (c *cache) getStruct(channel string, res interface{}) (bool, error) {
	return crawling.GetStruct(c.key(channel), res)
}

func (c *cache) getString(channel string) (string, bool) {
	result, err := crawling.GetString(c.key(channel))
	if err != nil {
		logger.Errorf("從 redis 獲取 %s 時出現錯誤: %v, 將嘗試使用本地快取", c.prefix, err)
		if str, ok := c.local[channel]; ok {
			return str.(string), ok
		}
	}
	return result, result != ""
}

func (c *cache) setString(channel, value string) error {
	err := crawling.SetString(c.key(channel), value)
	if err != nil {
		logger.Errorf("儲存 %s 到 redis 時出現錯誤: %v, 將嘗試使用本地快取", c.prefix, err)
		c.local[channel] = value
	}
	return err
}

func (c *cache) GetAsLiveStatus(channel string) (*ChannelStatus, bool) {
	liveStatus := &ChannelStatus{}
	ok, err := c.getStruct(channel, liveStatus)
	if err != nil {
		logger.Errorf("從 redis 獲取 %s 時出現錯誤: %v, 將嘗試使用本地快取", c.prefix, err)
		if status, ok := c.local[channel]; ok {
			return status.(*ChannelStatus), ok
		}
	}
	return liveStatus, ok
}
