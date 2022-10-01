package valorant

import (
	"context"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/logging"
	"sync"
)

const Tag = "valorant"

var (
	logger   = logging.GetCrawlerLogger(Tag)
	instance = &crawler{}
)

type crawler struct {
}

func (c *crawler) Listen(room string, publish crawling.Publisher, wg *sync.WaitGroup) context.CancelFunc {
	//TODO implement me
	panic("implement me")
}

func (c *crawler) Prefix() string {
	//TODO implement me
	panic("implement me")
}

func (c *crawler) IsValidTopic(topic string) bool {
	//TODO implement me
	panic("implement me")
}

func (c *crawler) Init() {
	//TODO implement me
	panic("implement me")
}

func (c *crawler) Start() {
	//TODO implement me
	panic("implement me")
}

func (c *crawler) Stop(wg *sync.WaitGroup) {
	//TODO implement me
	panic("implement me")
}
