package bilibili

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWebSocket(t *testing.T) {
	// no need load because default settings is valid
	// file.LoadYaml("bilibili", bilibiliYaml)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	wg := &sync.WaitGroup{}
	go startWebSocket(ctx, wg)

	listen := []string{
		"876396",
	}

	go subscribeAll(listen, ctx, cancel, nil)
	<-time.After(time.Second * 30)
	cancel()
	wg.Wait()
}
