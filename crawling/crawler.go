package crawling

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

type Publisher func(room string, arg interface{})

type Crawler interface {
	Prefix() string

	IsValidTopic(topic string) bool

	// Init 加載 Yaml
	Init()

	// Start 加載 Spider
	Start()

	Listen(room string, publish Publisher) context.CancelFunc

	// Stop 關閉所有 Spider
	Stop(wg *sync.WaitGroup)
}

type Crawling struct {
	Crawler   Crawler
	Name      string
	Logger    *logrus.Entry
	Listening map[string]context.CancelFunc
}

var crawlers = make(map[string]*Crawling)

func RegisterCrawler(name string, crawler Crawler, logger *logrus.Entry) {
	crawlers[name] = &Crawling{
		Crawler:   crawler,
		Name:      name,
		Logger:    logger,
		Listening: make(map[string]context.CancelFunc),
	}
}

func GetCrawler(name string) *Crawling {
	if craw, ok := crawlers[name]; ok {
		return craw
	} else {
		return nil
	}
}

func InitAllCrawlers() {
	for _, crawling := range crawlers {
		crawling.Crawler.Init()
	}
}
