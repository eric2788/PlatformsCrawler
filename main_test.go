package main

import (
	"context"
	"github.com/eric2788/PlatformsCrawler/crawlers/bilibili"
	"github.com/eric2788/PlatformsCrawler/crawlers/twitter"
	"github.com/eric2788/PlatformsCrawler/crawlers/youtube"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"testing"
	"time"

	_ "net/http/pprof"
)

func TestPprof(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:45678", http.DefaultServeMux))
	}()
	main()
}

func TestDisableCrawlers(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)

	file.LoadApplicationYaml()

	file.ApplicationYaml.DisabledCrawlers = []string{twitter.Tag, bilibili.Tag, youtube.Tag}

	crawling.InitAllCrawlers()

	ticker := time.NewTicker(time.Second * time.Duration(file.ApplicationYaml.CheckInterval))
	ctx, cancel := context.WithCancel(context.Background())

	waitStop := make(chan struct{}, 1)
	go crawling.StartCrawling(ticker, ctx, waitStop)

	<-time.After(time.Second * 15)
	cancel()
	ticker.Stop()
	<-waitStop
}
