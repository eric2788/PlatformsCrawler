package twitter

import (
	"context"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/eric2788/PlatformsCrawler/logging"
	"sync"
	twitter "github.com/n0madic/twitter-scraper"
)

const Tag = "twitter"

var (
	logger   = logging.GetCrawlerLogger(Tag)
	instance = &crawler{}
	scraper = twitter.New()
)

type crawler struct {
}

func (c *crawler) Prefix() string {
	return "twitter"
}

func (c *crawler) IsValidTopic(topic string) bool {
	return true
}

func (c *crawler) Init() {
	file.LoadYaml("twitter", twitterYaml)
}

func (c *crawler) Start() {
	err := scraper.Login(twitterYaml.Username, twitterYaml.Password)
	if err != nil {
		logger.Errorf("使用用戶名 %s 登入推特失敗: %v, 將改用匿名登入", twitterYaml.Username, err)
		err = scraper.LoginOpenAccount()
	}
	if err != nil {
		logger.Errorf("爬蟲初始化失敗，本爬蟲可能無法正常運作")
		return
	}
	logger.Infof("Twitter 爬蟲已啟動")
}

func (c *crawler) Listen(username string, publish crawling.Publisher, wg *sync.WaitGroup) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go listenUserTweets(ctx, username, wg, publish)
	return cancel
}


func (c *crawler) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Infof("Twitter 爬蟲已關閉")
}

func init() {
	crawling.RegisterCrawler(Tag, instance, logger)
}
