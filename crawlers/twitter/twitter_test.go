package twitter

import (
	"context"
	"fmt"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"testing"
	"time"
)

func TestTwitterStream(t *testing.T) {

	file.LoadYaml("twitter", twitterYaml)

	crawling.InitRedis()
	startTwitterClient()

	initMuxHandle(func(room string, arg interface{}) {
		fmt.Printf("Send: twitter:%s => %+v\n", room, arg)
	})

	go refreshTwitterStream([]string{"Every3Minutes"})
	ctx, cancel := context.WithCancel(context.Background())
	waitStop, done := context.WithCancel(context.Background())
	go signalForStop(ctx, done)
	<-time.After(time.Minute * 10)
	cancel()
	<-waitStop.Done()
}
