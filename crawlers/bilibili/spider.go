package bilibili

import (
	"context"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const id = "platforms_crawler"

var (
	publisher crawling.Publisher
	livedSet  = mapset.NewSet()
	listening []string
)

// antiDuplicateLive 基於 LIVE 指令可能會連續發送幾次
func antiDuplicateLive(roomId float64) {
	livedSet.Add(roomId)
	<-time.After(time.Minute * time.Duration(bilibiliYaml.AntiDuplicateLive))
	livedSet.Remove(roomId)
}

// handleMessage here to publish redis message
func handleMessage(b []byte) {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		logger.Warnf("解析 JSON 數據時出現錯誤: %v", err)
		return
	}
	if info, ok := data["live_info"].(map[string]interface{}); ok {
		roomId := info["room_id"].(float64)
		// 有機會為 null
		if publisher != nil {

			if data["command"] == "LIVE" {
				if m, ok := data["content"].(map[string]interface{}); ok {
					if _, exist := m["live_time"]; !exist {
						return // 沒有 live_time 的 key 為多餘的開播推送
					}
				} else {
					logger.Warnf("無法把 content 轉換為 map (空 JSON 內容?), 已使用內置去重方式。")
					// 保險起見的方式
					if livedSet.Contains(roomId) {
						return // 開播通知去重
					} else {
						go antiDuplicateLive(roomId)
					}
				}
			}

			publisher(fmt.Sprintf("%d", int64(roomId)), b)
		} else {
			logger.Debugf("推送方式為 null，已略過")
		}

		// 僅作為 logging
		if data["command"] == "LIVE" {
			logger.Infof("檢測到 %s(%d) 在 B站 開播了。", info["name"], int64(roomId))
		} else {
			logger.Debugf("Received %s command from room %d", data["command"], int64(roomId))
		}
	} else {
		logger.Warnf("未知的房間 %+v", data["live_info"])
	}
}

func doSubscribeRequest(room []string) (url.URL, error) {

	httpUrl := url.URL{
		Host:     bilibiliYaml.BiliLiveHost,
		Path:     "/subscribe",
		RawQuery: "validate=false",
	}

	if bilibiliYaml.UseTLS {
		httpUrl.Scheme = "https"
	} else {
		httpUrl.Scheme = "http"
	}

	form := url.Values{
		"subscribes": room,
	}

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequest(http.MethodPost, httpUrl.String(), body)

	if err != nil {
		return httpUrl, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", id)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return httpUrl, err
	}

	if resp.StatusCode != 200 {
		return httpUrl, fmt.Errorf(resp.Status)
	}
	return httpUrl, nil
}

func subscribeAll(room []string, ctx context.Context, done context.CancelFunc, p crawling.Publisher) {

	if publisher == nil {
		publisher = p
	}

	retry := func() {
		logger.Warnf("三十秒後嘗試")
		select {
		case <-time.After(time.Second * 30):
			subscribeAll(room, ctx, done, p)
		case <-ctx.Done(): // 等待三十秒時需要刷新訂閱，則直接關閉
			done()
		}
	}

	logger.Debugf("正在設置訂閱...")
	httpUrl, err := doSubscribeRequest(room)
	listening = room

	if err != nil {
		logger.Errorf("嘗試設置訂閱時出現錯誤: %v", err)
		retry()
		return
	}

	logger.Debugf("設置訂閱成功。")

	defer done()
	<-ctx.Done()
	unSubscribe(httpUrl)
}

func unSubscribe(httpUrl url.URL) {

	logger.Debugf("正在清除訂閱...")

	req, err := http.NewRequest(http.MethodDelete, httpUrl.String(), nil)

	if err != nil {
		logger.Errorf("請求刪除先前的訂閱時出現錯誤: %v", err)
	}

	_, err = http.DefaultClient.Do(req)

	if err != nil {
		logger.Errorf("刪除先前的訂閱時出現錯誤: %v", err)
	}

	logger.Debugf("清除訂閱成功。")
}
