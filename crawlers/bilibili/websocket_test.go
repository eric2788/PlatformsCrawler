package bilibili

import (
	"context"
	"testing"
	"time"
)

func TestWebSocket(t *testing.T) {
	// no need load because default settings is valid
	// file.LoadYaml("bilibili", bilibiliYaml)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	go startWebSocket(ctx)

	listen := []string{
		"876396",
	}

	go subscribeAll(listen, ctx, cancel, nil)
}
