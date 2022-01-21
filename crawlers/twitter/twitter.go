package twitter

import (
	"context"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/eric2788/PlatformsCrawler/logging"
	"sync"
)

const Tag = "twitter"

var (
	logger   = logging.GetCrawlerLogger(Tag)
	instance = &crawler{}
)

type crawler struct {
}

func (c *crawler) Prefix() string {
	return "twitter"
}

func (c *crawler) IsValidTopic(topic string) bool {
	return true
}

func (c *crawler) Init() {
	file.LoadYaml("twitter", twitterYaml)
}

func (c *crawler) Start() {
	startTwitterClient()
	logger.Infof("Twitter 爬蟲已啟動")
}

func (c *crawler) ListenAll(room []string, publisher crawling.Publisher, done context.CancelFunc) context.CancelFunc {
	if mux == nil {
		initMuxHandle(publisher)
	}
	go refreshTwitterStream(room)
	ctx, cancel := context.WithCancel(context.Background())
	go signalForStop(ctx, done)
	return cancel
}

func (c *crawler) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Infof("Twitter 爬蟲已關閉")
}

func init() {
	crawling.RegisterCrawler(Tag, instance, logger)
}
