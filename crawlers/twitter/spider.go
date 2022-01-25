package twitter

import (
	"context"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/common-utils/set"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"time"
)

var (
	cbg       = context.Background()
	mux       *twitter.SwitchDemux
	client    *twitter.Client
	stream    *twitter.Stream
	listening *set.StringSet
)

// oauth1Client user auth
func oauth1Client() *http.Client {
	token := oauth1.NewToken(twitterYaml.AccessToken, twitterYaml.AccessTokenSecret)
	config := oauth1.Config{
		ConsumerKey:    twitterYaml.ConsumerKey,
		ConsumerSecret: twitterYaml.ConsumerSecret,
	}
	return config.Client(cbg, token)
}

// oauth2Client app auth
func oauth2Client() *http.Client {
	config := &clientcredentials.Config{
		ClientID:     twitterYaml.ConsumerKey,
		ClientSecret: twitterYaml.ConsumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	return config.Client(cbg)
}

func startTwitterClient() {
	client = twitter.NewClient(oauth1Client())
}

func refreshTwitterStream(screenNames []string) {

	userMap, err := UserLookUpCache(screenNames)

	if err != nil {
		logger.Warnf("查找用戶 %v 的ID時出現錯誤: %v", screenNames, err)
		logger.Warnf("十秒後重試...")
		<-time.After(time.Second * 10)
		refreshTwitterStream(screenNames)
		return
	}

	toFollow := make([]string, 0)

	for _, id := range userMap {
		toFollow = append(toFollow, id)
	}

	logger.Infof("正在刷新推特串流...")

	onlyFollowers := &twitter.StreamFilterParams{
		//FilterLevel:   "low",
		Follow:        toFollow,
		StallWarnings: twitter.Bool(true),
	}

	stream, err = client.Streams.Filter(onlyFollowers)

	if err != nil {
		logger.Warnf("嘗試刷新串流時出現錯誤: %v", err)
		logger.Warnf("十秒後重試...")
		<-time.After(time.Second * 10)
		refreshTwitterStream(screenNames)
		return
	}

	go mux.HandleChan(stream.Messages)

	logger.Infof("新的推特串流已啟動。")

}

func initMuxHandle(publisher crawling.Publisher) {
	m := twitter.NewSwitchDemux()
	m.StatusDeletion = func(deletion *twitter.StatusDeletion) {
		logger.Debugf("%s 刪除了動態", deletion.UserIDStr)
	}
	m.Tweet = func(tweet *twitter.Tweet) {
		if !listening.Contains(tweet.User.ScreenName) {
			return
		}
		logger.Infof("%s 發佈了新動態", tweet.User.Name)
		go publisher(tweet.User.ScreenName, tweet)
	}
	m.Warning = func(warning *twitter.StallWarning) {
		logger.Warnf("收到警告: %v", warning.Message)
	}
	mux = &m
}

func signalForStop(ctx context.Context, done context.CancelFunc) {
	defer done()
	<-ctx.Done()
	if stream != nil {
		logger.Infof("正在關閉推特串流...")
		stream.Stop() // blocking wait
		logger.Infof("推特串流已關閉。")
	}
}
