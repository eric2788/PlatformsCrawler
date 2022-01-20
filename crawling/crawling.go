package crawling

import (
	"PlatformsCrawler/config"
	"PlatformsCrawler/logging"
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/go-redis/redis/v8"
	"strings"
	"sync"
	"time"
)

var (
	logger          = logging.GetMainLogger()
	ctx             = context.Background()
	exceptionTopics = mapset.NewSet()
)

// StartCrawling remember use via go
func StartCrawling(tick *time.Ticker, ctx context.Context) {

	rConfig := config.ApplicationYaml.Redis

	rcli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rConfig.Host, rConfig.Port),
		Password: rConfig.Password,
		DB:       rConfig.Database,
	})

	defer logger.Infof("所有爬蟲已關閉。")
	defer stopAllAndWait()

	// 啟動所有爬蟲
	for _, crawling := range crawlers {
		crawling.Crawler.Start()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-tick.C:
			if !ok {
				return
			}
			for _, crawling := range crawlers {
				publisher := getPublisherFunc(rcli, crawling)
				go crawlEach(rcli, crawling, publisher)
			}
		}
	}
}

func getPublisherFunc(client *redis.Client, crawling *Crawling) Publisher {
	return func(room string, arg interface{}) {
		prefix := crawling.Crawler.Prefix()
		topic := fmt.Sprintf("%s:%s", prefix, room)
		if err := client.Publish(ctx, topic, arg).Err(); err != nil {
			logger.Errorf("嘗試推送訊息到 topic %s 時出現錯誤: %v", topic, err)
			logger.Errorf("推送訊息: %+v", arg)

			// does it need cancel the pubsub when error ?
		}
	}
}

func crawlEach(client *redis.Client, crawling *Crawling, publisher Publisher) {

	listening, crawler, cLogger := crawling.Listening, crawling.Crawler, crawling.Logger

	topics, err := client.PubSubChannels(ctx, fmt.Sprintf("%s:*", crawler.Prefix())).Result()

	if err != nil {
		cLogger.Warnf("嘗試獲取 %s 的 pubsub channels 時出現錯誤: %v", crawler.Prefix(), err)
		return
	}

	topicSet := mapset.NewSet()

	for _, channel := range topics {
		topic := strings.ReplaceAll(channel, fmt.Sprintf("%s:", crawler.Prefix()), "")
		if exceptionTopics.Contains(topic) {
			continue
		} else if !crawler.IsValidTopic(topic) {
			cLogger.Warnf("%s 不是一個有效的 topic, 已略過。", topic)
			exceptionTopics.Add(topic)
			continue
		} else {
			topicSet.Add(topic)
		}
	}

	listeningSet := mapset.NewSet()

	for key := range listening {
		listeningSet.Add(key)
	}

	toListen := topicSet.Difference(listeningSet)
	stopListen := listeningSet.Difference(topicSet)

	for listen := range toListen.Iter() {
		t := listen.(string)
		canceller := crawler.Listen(t, publisher)
		listening[t] = canceller
	}

	for stop := range stopListen.Iter() {
		t := stop.(string)
		if cancel, ok := listening[t]; ok {
			cancel()
			delete(listening, t)
		} else {
			cLogger.Errorf("嘗試停止 topic %s 時 發現 不存在於監聽列表。", t)
		}
	}

}

func stopAllAndWait() {
	gp := &sync.WaitGroup{}
	gp.Add(len(crawlers))
	for _, crawling := range crawlers {
		go crawling.Crawler.Stop(gp)
	}
	gp.Wait()
}
