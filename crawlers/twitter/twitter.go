package twitter

import (
	"context"
	"fmt"
	"sync"

	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/eric2788/PlatformsCrawler/logging"
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
	scraper.WithReplies(true)
	scraper.WithDelay(twitterYaml.RequestDelay)
}

func (c *crawler) Start() {
	err := scraper.Login(twitterYaml.Username, twitterYaml.Password, twitterYaml.EmailCode)
	if err != nil {
		logger.Errorf("使用用戶名 %s 登入推特失敗: %v, 將改用匿名登入", twitterYaml.Username, err)
		err = scraper.LoginOpenAccount()
	}
	if err != nil {
		logger.Errorf("爬蟲初始化失敗，本爬蟲可能無法正常運作: %v", err)
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
	crawling.AddCommandHandler("twitter-login", func(command crawling.CommandSchema) string {
		if len(command.Args) != 3 {
			return "參數錯誤，請輸入用戶名、密碼、Code/Email"
		}
		args := command.Args
		username, password, code := args[0], args[1], args[2]
		err := scraper.Login(username, password, code)
		if err != nil {
			return fmt.Sprintf("登入失敗: %v", err)
		}else {
			return "登入成功"
		}
	})
}