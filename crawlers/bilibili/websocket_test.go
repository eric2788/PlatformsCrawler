package bilibili

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
	"time"
)

var testListen = []string{
	"876396",
	"22571958",
	"21320551",
}

func TestWebSocket(t *testing.T) {
	// no need load because default settings is valid
	// file.LoadYaml("bilibili", bilibiliYaml)

	logrus.SetLevel(logrus.DebugLevel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	wg := &sync.WaitGroup{}
	go startWebSocket(ctx, wg)

	go subscribeAll(testListen, ctx, cancel, nil)
	<-time.After(time.Second * 30)
	cancel()
	wg.Wait()
}

func TestReuseRequest(t *testing.T) {
	for i := 0; i < 3; i++ {
		if _, err := doSubscribeRequest(testListen); err != nil {
			t.Fatal(err)
		}
	}
}
