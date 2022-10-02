package valorant

import (
	"context"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/eric2788/PlatformsCrawler/logging"
	"github.com/google/uuid"
	"sync"
)

const Tag = "valorant"

var (
	logger   = logging.GetCrawlerLogger(Tag)
	instance = &crawler{}
)

type crawler struct {
}

func (c *crawler) Listen(uuid string, publish crawling.Publisher, wg *sync.WaitGroup) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go runValorantMatchTrack(ctx, uuid, wg, publish)
	return cancel
}

func (c *crawler) Prefix() string {
	return "valorant"
}

func (c *crawler) IsValidTopic(topic string) bool {
	_, err := uuid.Parse(topic)
	return err == nil
}

func (c *crawler) Init() {
	file.LoadYaml("valorant", valorantYaml)
}

func (c *crawler) Start() {
	logger.Infof("Valorant 爬蟲已啟動。")
}

func (c *crawler) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Infof("Valorant 爬蟲已關閉")
}

func init() {
	crawling.RegisterCrawler(Tag, instance, logger)
}
