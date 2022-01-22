package bilibili

import (
	"context"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
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
	stop context.CancelFunc
	wg   *sync.WaitGroup
}

func (c *crawler) Prefix() string {
	return "blive"
}

func (c *crawler) IsValidTopic(topic string) bool {
	_, err := strconv.Atoi(topic)
	return err == nil
}

func (c *crawler) Init() {
	file.LoadYaml("bilibili", bilibiliYaml)
}

func (c *crawler) Start() {
	ctx, stop := context.WithCancel(context.Background())
	c.wg = &sync.WaitGroup{}
	go startWebSocket(ctx, c.wg)
	c.stop = stop
	logger.Infof("Bilibili 爬蟲已啟動")
}

func (c *crawler) ListenAll(room []string, publisher crawling.Publisher, done context.CancelFunc) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go subscribeAll(room, ctx, done, publisher)
	return cancel
}

func (c *crawler) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	c.stop()
	c.wg.Wait()
	logger.Infof("Bilibili 爬蟲已關閉")
}

func init() {
	crawling.RegisterCrawler(Tag, instance, logger)
}
