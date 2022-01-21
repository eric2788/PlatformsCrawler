package crawling

import (
	"context"
	mapset "github.com/deckarep/golang-set"
	"github.com/sirupsen/logrus"
	"sync"
)

// Publisher functional interface
type Publisher func(room string, arg interface{})

// interface
type (
	EachCrawling interface {

		// Listen 逐一監聽, 與 ListenAll 二選一
		Listen(room string, publish Publisher, wg *sync.WaitGroup) context.CancelFunc
	}

	OnceCrawling interface {

		// ListenAll 一次性監聽, 與 Listen 二選一
		ListenAll(room []string, publisher Publisher, done context.CancelFunc) context.CancelFunc
	}
)

// struct
type (
	EachSpec struct {
		listening map[string]context.CancelFunc
		wg        *sync.WaitGroup
	}

	OnceSpec struct {
		waitStop context.Context
		stopAll  context.CancelFunc
	}
)

type Crawler interface {
	Prefix() string

	IsValidTopic(topic string) bool

	// Init 加載 Yaml
	Init()

	// Start 加載 Spider
	Start()

	// Stop 關閉所有 Spider
	Stop(wg *sync.WaitGroup)
}

type Crawling struct {
	Crawler   Crawler
	Name      string
	Logger    *logrus.Entry
	Listening mapset.Set
	spec      interface{}
}

var crawlers = make(map[string]*Crawling)

func RegisterCrawler(name string, crawler Crawler, logger *logrus.Entry) {

	var spec interface{}

	switch crawler.(type) {
	case EachCrawling:
		spec = &EachSpec{
			listening: make(map[string]context.CancelFunc),
			wg:        &sync.WaitGroup{},
		}
	case OnceCrawling:
		spec = &OnceSpec{}
	default:
		panic("unknown crawling type")
	}

	crawlers[name] = &Crawling{
		Crawler:   crawler,
		Name:      name,
		Logger:    logger,
		Listening: mapset.NewSet(),
		spec:      spec,
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
