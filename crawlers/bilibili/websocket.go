package bilibili

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

func startWebSocket(ctx context.Context, wg *sync.WaitGroup) {

	wsUrl := url.URL{
		Host:     bilibiliYaml.BiliLiveHost,
		Path:     "/ws",
		RawQuery: fmt.Sprintf("id=%s", id),
	}

	if bilibiliYaml.UseTLS {
		wsUrl.Scheme = "wss"
	} else {
		wsUrl.Scheme = "ws"
	}

	con, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)

	if err != nil {
		logger.Errorf("連線到 Websocket %s 時出現錯誤: %v", wsUrl.String(), err)
		logger.Warnf("十秒後重試")
		<-time.After(time.Second * 10)
		startWebSocket(ctx, wg)
		return
	}

	con.SetCloseHandler(func(code int, text string) error {
		return con.WriteMessage(websocket.CloseMessage, nil)
	})

	wg.Add(1)
	onReceiveMessage(ctx, con, wg)
}

func onReceiveMessage(ctx context.Context, conn *websocket.Conn, wg *sync.WaitGroup) {
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf("關閉 Websocket 時出現錯誤: %v", err)
		} else {
			logger.Debugf("連接關閉成功。")
		}
		wg.Done()
	}()
	logger.Infof("Biligo WebSocket 連接成功。")
	for {
		select {
		case <-ctx.Done():
			logger.Infof("正在關閉 WebSocket...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "stop"))
			if err != nil {
				logger.Errorf("發送 websocket 關閉訊息時出現錯誤: %v", err)
			}
			return
		default:
			_, message, err := conn.ReadMessage()
			// Error
			if err != nil {
				logger.Errorf("Websocket 嘗試讀取消息時出現錯誤: %v", err)
				go retryDelay(ctx, wg)
				return
			}
			go handleMessage(message)
		}
	}
}

// test only
func testReadMessageWithRandomError(conn *websocket.Conn) (messageType int, p []byte, err error) {
	messageType, p, err = conn.ReadMessage()
	rand.Seed(time.Now().UnixMicro())
	if rand.Intn(20) == 4 {
		err = fmt.Errorf("test error")
	}
	return
}

func retryDelay(ctx context.Context, wg *sync.WaitGroup) {
	logger.Warnf("五秒後重連...")
	<-time.After(time.Second * 5)
	startWebSocket(ctx, wg)
	if listening != nil {
		// 重新訂閱
		for _, err := doSubscribeRequest(listening); err != nil; _, err = doSubscribeRequest(listening) {
			logger.Errorf("重新訂閱失敗: %v，五秒後重試...", err)
			<-time.After(time.Second * 5)
		}
		logger.Infof("重新訂閱成功。")
	}
}
