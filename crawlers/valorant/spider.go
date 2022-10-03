package valorant

import (
	"context"
	"fmt"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/common-utils/datetime"
	"sync"
	"time"
)

func saveLatestMatch(uuid, matchId string) {
	key := fmt.Sprintf("valorant:latest_match:%s", uuid)
	err := crawling.SetString(key, matchId)
	if err != nil {
		logger.Errorf("嘗試保存玩家 %s 的最新對戰記錄時出現錯誤: %v", uuid, err)
	}
}

func getLatestMatch(uuid string) (string, bool) {
	key := fmt.Sprintf("valorant:latest_match:%s", uuid)
	latest, err := crawling.GetString(key)

	if err != nil || latest == "" {

		if err != nil {
			logger.Errorf("嘗試獲取玩家 %s 的最新對戰記錄時出現錯誤: %v", uuid, err)
		}

		return "", false
	} else {
		return latest, true
	}
}

func getDisplayNameByUuid(uuid string) (string, bool) {
	key := fmt.Sprintf("valorant:display_name:%s", uuid)
	displayName, err := crawling.GetString(key)

	// redis 快取找到
	if displayName != "" && err == nil {
		return displayName, true
	}

	if err != nil {
		logger.Errorf("嘗試獲取玩家 %s 的顯示名稱時出現錯誤: %v", uuid, err)
	} else if displayName == "" {
		logger.Warnf("玩家 %s 的顯示名稱不在快取中或已過期。", uuid)
	}

	logger.Warnf("將使用 API 請求獲取 %s 的 顯示名稱。", uuid)

	account, err := getDisplayName(uuid)
	if err != nil {
		logger.Errorf("嘗試獲取玩家 %s 的顯示名稱時出現錯誤: %v", uuid, err)
		return "", false
	} else {
		displayName = fmt.Sprintf("%s#%s", account.Name, account.Tag)
		err = crawling.SetStringTemp(key, displayName, time.Hour*24*20)
		if err != nil {
			logger.Errorf("嘗試保存玩家 %s 的顯示名稱大到redis時出現錯誤: %v", uuid, err)
		}
		return displayName, true
	}
}

func runValorantMatchTrack(ctx context.Context, uuid string, wg *sync.WaitGroup, publish crawling.Publisher) {

	ticker := time.NewTicker(time.Duration(valorantYaml.Interval) * time.Second)
	displayName, success := getDisplayNameByUuid(uuid)
	if !success {
		displayName = fmt.Sprintf("(%s)", uuid)
	}

	defer wg.Done()
	defer logger.Infof("玩家 %s 監控中止。", displayName)
	defer ticker.Stop()

	logger.Infof("玩家 %s 的 Valorant 遊戲狀態監控已啟動。", displayName)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			matches, err := getValorantMatches(uuid)
			if err != nil {
				logger.Errorf("嘗試獲取玩家 %s 的遊戲狀態時出現錯誤: %v", displayName, err)
				continue
			} else if len(matches) == 0 {
				logger.Warnf("玩家 %s 的對戰記錄為空，已略過。", displayName)
			}

			latestData := matches[0].MetaData
			logger.Debugf("玩家 %s 的最新對戰資料ID為 %v, 時間: %s", displayName, latestData.MatchId, datetime.FormatSeconds(latestData.GameStart))

			lastMatchId, ok := getLatestMatch(uuid)

			// 與上一次的狀態相同 => 忽略
			if ok && lastMatchId == latestData.MatchId {
				continue
			}

			saveLatestMatch(uuid, latestData.MatchId)

			// 尚未有上一次檢測的資料 + 最新對戰記錄距今已超過1小時 => 忽略
			if !ok && datetime.Duration(latestData.GameStart, time.Now().Unix()).Hours() > 1 {
				continue
			}

			// 從快取獲取顯示名稱，如果快取過期則可獲取最新的顯示名稱
			latestDisplayName, latestSuccess := getDisplayNameByUuid(uuid)

			// 如果不成功，則使用回第一次索取的顯示名稱
			if !latestSuccess {
				latestDisplayName = displayName
			}

			logger.Infof("玩家 %s 的最新對戰訊息已更新。最新時間為: %s", latestDisplayName, datetime.FormatSeconds(latestData.GameStart))

			publishData := &MatchMetaDataPublish{
				Data:        &latestData,
				DisplayName: latestDisplayName,
			}

			publish(uuid, publishData)
		}
	}

}
