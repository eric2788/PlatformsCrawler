package main

import (
	"context"
	"flag"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/eric2788/PlatformsCrawler/rest"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"

	_ "github.com/eric2788/PlatformsCrawler/crawlers/bilibili"
	_ "github.com/eric2788/PlatformsCrawler/crawlers/twitter"
	_ "github.com/eric2788/PlatformsCrawler/crawlers/youtube"
)

var debug = flag.Bool("debug", os.Getenv("DEBUG") == "true", "enable debug level")
var port = flag.Int("port", 8989, "the restful api port")

func main() {

	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	file.LoadApplicationYaml()

	crawling.InitAllCrawlers()

	ticker := time.NewTicker(time.Second * time.Duration(file.ApplicationYaml.CheckInterval))
	ctx, cancel := context.WithCancel(context.Background())

	waitStop := make(chan struct{}, 1)
	go crawling.StartCrawling(ticker, ctx, waitStop)
	go rest.StartServe(*port)
	go debugServe()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	cancel()
	ticker.Stop()
	<-waitStop
}
