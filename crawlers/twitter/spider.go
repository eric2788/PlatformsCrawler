package twitter

import (
	"context"
	"sync"
	"time"

	"github.com/eric2788/PlatformsCrawler/crawling"
)


func listenUserTweets(ctx context.Context, username string, wg *sync.WaitGroup, publisher crawling.Publisher){
	
	ticker := time.NewTicker(time.Duration(twitterYaml.ScrapeInterval) * time.Second)
	
	defer wg.Done()
	defer logger.Infof("推特 %s 監控中止。", username)
	defer ticker.Stop()

	logger.Infof("推特 %s 監控已啟動。", username)

	for {
		select{
		case <-ctx.Done():
			return
		case <-ticker.C:
			tweetsChan := scraper.GetTweets(ctx, username, 1)
			lastTweet, ok := <-tweetsChan 
			if !ok {
				logger.Warnf("找不到用戶 %s 的推文。", username)
				continue
			}
			if lastTweet.Error != nil {
				logger.Errorf("刷取用戶 %s 推文內容時出現錯誤 %v", username, lastTweet.Error)
				continue
			}
			go publisher(username, lastTweet)
		}
	}
}