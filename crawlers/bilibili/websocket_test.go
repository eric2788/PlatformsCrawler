package bilibili

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
	"time"
)

func TestWebSocket(t *testing.T) {
	// no need load because default settings is valid
	// file.LoadYaml("bilibili", bilibiliYaml)

	logrus.SetLevel(logrus.DebugLevel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	wg := &sync.WaitGroup{}
	go startWebSocket(ctx, wg)

	listen := []string{
		"876396",
		"22571958",
		"21320551",
	}

	go subscribeAll(listen, ctx, cancel, nil)
	<-time.After(time.Second * 30)
	cancel()
	wg.Wait()
}
