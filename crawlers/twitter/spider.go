package twitter

import (
	"context"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"time"
)

var (
	cbg    = context.Background()
	mux    *twitter.SwitchDemux
	client *twitter.Client
	stream *twitter.Stream
)

func oauth1Client() *http.Client {
	token := oauth1.NewToken(twitterYaml.AccessToken, twitterYaml.AccessTokenSecret)
	config := oauth1.Config{
		ConsumerKey:    twitterYaml.ConsumerKey,
		ConsumerSecret: twitterYaml.ConsumerSecret,
	}
	return config.Client(cbg, token)
}

func oauth2Client() *http.Client {
	config := &clientcredentials.Config{
		ClientID:     twitterYaml.ConsumerKey,
		ClientSecret: twitterYaml.ConsumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	return config.Client(cbg)
}

func startTwitterClient() {
	client = twitter.NewClient(oauth2Client())
}

func refreshTwitterStream(screenNames []string) {

	// 如果之前已有串流
	if stream != nil {
		stream.Stop()
		stream = nil
	}

	userMap, err := UserLookUpCache(screenNames)

	if err != nil {
		logger.Warnf("查找用戶ID時出現錯誤: %v", err)
		logger.Warnf("十秒後重試...")
		<-time.After(time.Second * 10)
		refreshTwitterStream(screenNames)
		return
	}

	toFollow := make([]string, 0)

	for _, id := range userMap {
		toFollow = append(toFollow, id)
	}

	onlyFollowers := &twitter.StreamFilterParams{
		FilterLevel:   "medium",
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

	mux.HandleChan(stream.Messages)
}

func initMuxHandle(publisher crawling.Publisher) {
	m := twitter.NewSwitchDemux()
	m.StatusDeletion = func(deletion *twitter.StatusDeletion) {
		logger.Infof("%s 刪除了動態", deletion.UserIDStr)
	}
	m.Tweet = func(tweet *twitter.Tweet) {
		logger.Infof("%s 發佈了新動態", tweet.User.Name)
		logger.Infof("%+v", *tweet)
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
		stream.Stop() // blocking wait
	}
}
