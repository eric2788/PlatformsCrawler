package main

import (
	"PlatformsCrawler/config"
	"PlatformsCrawler/crawling"
	"context"
	"os"
	"os/signal"
	"time"
)

func main() {

	config.LoadApplicationYaml()

	crawling.InitAllCrawlers()

	ticker := time.NewTicker(time.Second * time.Duration(config.ApplicationYaml.CheckInterval))
	ctx, cancel := context.WithCancel(context.Background())
	go crawling.StartCrawling(ticker, ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	cancel()
	ticker.Stop()
}
