package twitter

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/eric2788/PlatformsCrawler/crawling"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

var lastTweetIdCache = crawling.NewCache("twitter", "last_tweet_id")

var excepted = mapset.NewSet[string]()


type TweetContent struct {
	Tweet *twitterscraper.TweetResult `json:"tweet"`
	Profile *twitterscraper.Profile `json:"profile"`
	NickName string `json:"nick_name"`
}

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

			if excepted.Contains(username) {
				logger.Debugf("用戶 %s 已被列入例外，已跳過。", username)
				return
			}

			tweetsChan := scraper.GetTweets(ctx, username, 1)
			lastTweet, ok := <-tweetsChan 
			if !ok {
				logger.Warnf("找不到用戶 %s 的推文。", username)
				continue
			}
			if lastTweet.Error != nil {
				logger.Errorf("刷取用戶 %s 推文內容時出現錯誤 %v", username, lastTweet.Error)
				if strings.Contains(lastTweet.Error.Error(), "Not authorized to view") {
					logger.Warnf("用戶 %s 的推文不可見，將停止監控。", username)
					excepted.Add(username)
					return
				}
				continue
			}

			lastTweetIdFromCache, exist := lastTweetIdCache.GetString(username)

			if exist && lastTweetIdFromCache == lastTweet.ID {
				logger.Debugf("最新推文ID與上次發布的推文ID相同，已跳過。")
				continue
			}

			lastTweetIdCache.SetString(username, lastTweet.ID)

			profile, _ := getProfileByScreen(username)
			nickName := lastTweet.Name
			if nickName == "" {
				logger.Warnf("找不到用戶 %s 的顯示名稱, 將採用API請求", username)
				nick, exist := getDisplayNameByScreen(username)
				if !exist {
					logger.Warnf("找不到用戶 %s 的顯示名稱, 將用回 username", username)
					nick = username
				}
				nickName = nick
			}

			logger.Infof("%s 發佈了一則新推文", nickName)
			
			go publisher(username, &TweetContent{
				Tweet: lastTweet,
				NickName: nickName,
				Profile: profile,
			})
		}
	}
}

func getProfileByScreen(screen string) (*twitterscraper.Profile, bool) {
	key := fmt.Sprintf("twitter:profile:%s", screen)
	var profile = &twitterscraper.Profile{}
	exist, err := crawling.GetStruct(key, profile)

	// redis 快取找到
	if exist && err == nil {
		return profile, true
	}

	if err != nil {
		logger.Errorf("嘗試獲取玩家 %s 的個人檔案時出現錯誤: %v", screen, err)
	} else if profile == nil || !exist {
		logger.Warnf("玩家 %s 的個人檔案不在快取中或已過期。", screen)
	}

	logger.Warnf("將使用 API 請求獲取 %s 的 個人檔案。", screen)

	account, err := scraper.GetProfile(screen)
	if err != nil {
		logger.Errorf("嘗試獲取玩家 %s 的個人檔案時出現錯誤: %v", screen, err)
		return nil, false
	} else {
		profile = &account
		err = crawling.Store(key, profile)
		if err != nil {
			logger.Errorf("嘗試保存玩家 %s 的個人檔案到redis時出現錯誤: %v", screen, err)
		}
		return profile, true
	}
}


func getDisplayNameByScreen(screen string) (string, bool) {
	key := fmt.Sprintf("twitter:display_name:%s", screen)
	displayName, err := crawling.GetString(key)

	// redis 快取找到
	if displayName != "" && err == nil {
		return displayName, true
	}

	if err != nil {
		logger.Errorf("嘗試獲取玩家 %s 的顯示名稱時出現錯誤: %v", screen, err)
	} else if displayName == "" {
		logger.Warnf("玩家 %s 的顯示名稱不在快取中或已過期。", screen)
	}

	logger.Warnf("將使用 API 請求獲取 %s 的 顯示名稱。", screen)

	account, exist := getProfileByScreen(screen)
	if !exist {
		logger.Errorf("嘗試獲取玩家 %s 的顯示名稱時出現錯誤: 檔案不存在", screen)
		return "", false
	} else {
		displayName = account.Name
		err = crawling.SetStringTemp(key, displayName, time.Hour*24*30)
		if err != nil {
			logger.Errorf("嘗試保存玩家 %s 的顯示名稱到redis時出現錯誤: %v", screen, err)
		}
		return displayName, true
	}
}