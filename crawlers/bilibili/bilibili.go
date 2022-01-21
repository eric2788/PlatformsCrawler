package bilibili

import (
	"context"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/logging"
	"strconv"
	"sync"
)

const Tag = "bilibili"

var (
	logger   = logging.GetCrawlerLogger(Tag)
	instance = &crawler{}
)

type crawler struct {
}

func (c *crawler) Prefix() string {
	return "blive"
}

func (c *crawler) IsValidTopic(topic string) bool {
	_, err := strconv.Atoi(topic)
	return err == nil
}

func (c *crawler) Init() {

}

func (c *crawler) Start() {
	logger.Infof("Bilibili 爬蟲已啟動")
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
	logger.Infof("Bilibili 爬蟲已關閉")
}

func init() {
	crawling.RegisterCrawler(Tag, instance, logger)
}
