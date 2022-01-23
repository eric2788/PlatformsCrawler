package youtube

import (
	"context"
	mapset "github.com/deckarep/golang-set"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/eric2788/PlatformsCrawler/logging"
	"strings"
	"sync"
)

const Tag = "youtube"

var (
	logger   = logging.GetCrawlerLogger(Tag)
	instance = &crawler{
		channels:  mapset.NewSet(),
		nameCache: make(map[string]string),
	}
)

type crawler struct {
	channels  mapset.Set
	nameCache map[string]string
}

func (c *crawler) Prefix() string {
	return "ylive"
}

func (c *crawler) IsValidTopic(topic string) bool {
	return strings.HasPrefix(topic, "UC")
}

func (c *crawler) Init() {
	file.LoadYaml("youtube", youtubeYaml)
	initKeywordRegexp()
	initYoutubeService()
}

func (c *crawler) Start() {
	logger.Infof("Youtube 爬蟲已啟動")
}

func (c *crawler) Listen(room string, publish crawling.Publisher, wg *sync.WaitGroup) context.CancelFunc {
	logger.Infof("即將監控頻道 %s", room)
	ctx, cancel := context.WithCancel(context.Background())
	go runYoutubeSpider(ctx, room, wg, publish)
	c.channels.Add(room)
	return cancel
}

func (c *crawler) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Infof("Youtube 爬蟲已關閉")
}

func (c *crawler) getChannelName(id string) string {
	if name, ok := c.nameCache[id]; ok {
		return name
	} else {
		toFetch := make([]string, 0)
		for channel := range c.channels.Iter() {
			toFetch = append(toFetch, channel.(string))
		}
		channelNames, err := getChannelNames(toFetch...)

		if err != nil {
			logger.Errorf("嘗試刷新用戶頻道資訊時出現錯誤: %v", err)
			return id // 暫時返回頻道ID
		}

		for id, screen := range channelNames {
			c.nameCache[id] = screen
		}

		if name, ok = channelNames[id]; ok {
			return name
		} else {
			return id // 查無頻道
		}
	}
}

func init() {
	crawling.RegisterCrawler(Tag, instance, logger)
}
