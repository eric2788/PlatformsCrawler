package crawling

import (
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/eric2788/PlatformsCrawler/logging"
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
func StartCrawling(tick *time.Ticker, ctx context.Context, stop chan<- struct{}) {

	initRedis()

	defer stopAllAndWait(stop)

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
				publisher := getPublisherFunc(cli, crawling)
				go crawlEach(cli, crawling, publisher)
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

	listeningSet, spec, crawler, cLogger := crawling.Listening, crawling.spec, crawling.Crawler, crawling.Logger

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

	switch cc := crawler.(type) {
	// 每次新增監聽
	case EachCrawling:

		each := spec.(*EachSpec)
		listenerMap := each.listening
		waitGroup := each.wg

		toListen := topicSet.Difference(listeningSet)
		stopListen := listeningSet.Difference(topicSet)

		// 各自監控
		for listen := range toListen.Iter() {
			t := listen.(string)
			waitGroup.Add(1)
			canceller := cc.Listen(t, publisher, waitGroup)
			listenerMap[t] = canceller
			listeningSet.Add(t)
		}

		// 各自中止
		for stop := range stopListen.Iter() {
			t := stop.(string)
			if cancel, ok := listenerMap[t]; ok {
				cancel()
				delete(listenerMap, t)
				listeningSet.Remove(t)
			} else {
				cLogger.Errorf("嘗試停止 topic %s 時 發現 不存在於監聽列表。", t)
			}
		}
	// 一次性新增監聽
	case OnceCrawling:

		diff := listeningSet.SymmetricDifference(topicSet)

		i := 0
		for range diff.Iter() {
			i += 1
		}

		if i == 0 {
			return
		} else {
			cLogger.Infof("即將追加 %d 個監控", i)
		}

		once := spec.(*OnceSpec)

		// 先前已有啟動
		if once.stopAll != nil {
			// 先中止所有
			once.stopAll()
			// 等待先前的停止
			<-once.waitStop.Done()
		}

		// 再訂閱所有
		toListen := topicSet

		oneTimeListen := make([]string, 0)
		for listen := range toListen.Iter() {
			t := listen.(string)
			oneTimeListen = append(oneTimeListen, t)
		}

		runner, done := context.WithCancel(ctx)
		canceller := cc.ListenAll(oneTimeListen, publisher, done)

		crawling.Listening = topicSet

		once.stopAll = canceller
		once.waitStop = runner

	default:
		logger.Errorf("爬蟲 %s 沒有可用的監控方式。", crawling.Name)
	}

}

func stopCrawler(crawling *Crawling, wg *sync.WaitGroup) {

	spec := crawling.spec

	switch s := spec.(type) {
	case *EachSpec:
		for _, cancelFunc := range s.listening {
			cancelFunc()
		}
		s.wg.Wait()
	case *OnceSpec:
		s.stopAll()
		<-s.waitStop.Done()
	default:
		logger.Errorf("%s 沒有可用的關閉方式", crawling.Name)
	}

	wg.Done()
}

func stopAllAndWait(stop chan<- struct{}) {

	gp := &sync.WaitGroup{}

	// stop all topics for each crawler
	gp.Add(len(crawlers))
	for _, crawling := range crawlers {
		go stopCrawler(crawling, gp)
	}
	gp.Wait()

	// stop all crawlers
	gp.Add(len(crawlers))
	for _, crawling := range crawlers {
		go crawling.Crawler.Stop(gp)
	}
	gp.Wait()
	logger.Infof("所有爬蟲已關閉。")
	stop <- struct{}{}
}
