package youtube

import (
	"context"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/eric2788/PlatformsCrawler/logging"
	"strings"
	"sync"
)

const Tag = "youtube"

var (
	logger   = logging.GetCrawlerLogger(Tag)
	instance = &crawler{}
)

type crawler struct {
}

func (c *crawler) Prefix() string {
	return "ylive"
}

func (c *crawler) IsValidTopic(topic string) bool {
	return strings.HasPrefix(topic, "UC")
}

func (c *crawler) Init() {
	file.LoadYaml("youtube", youtubeYaml)
}

func (c *crawler) Start() {
	logger.Infof("Youtube 爬蟲已啟動")
}

func (c *crawler) Listen(room string, publish crawling.Publisher, wg *sync.WaitGroup) context.CancelFunc {
	logger.Infof("Listen %s", room)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		logger.Infof("Stop Listen %s", room)
		wg.Done()
	}()
	return cancel
}

func (c *crawler) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Infof("Youtube 爬蟲已關閉")
}

func init() {
	crawling.RegisterCrawler(Tag, instance, logger)
}
