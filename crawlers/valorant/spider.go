package valorant

import (
	"context"
	"fmt"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/common-utils/datetime"
	"sync"
	"time"
)

var (
	lastMatchMap = &sync.Map{}
)

func runValorantMatchTrack(ctx context.Context, name, tag string, wg *sync.WaitGroup, publish crawling.Publisher) {

	ticker := time.NewTicker(time.Duration(valorantYaml.Interval) * time.Second)
	displayName := fmt.Sprintf("%s#%s", name, tag)

	defer wg.Done()
	defer logger.Infof("玩家 %s 監控中止。", displayName)
	defer ticker.Stop()

	logger.Infof("玩家 %s 的 Valorant 遊戲狀態監控已啟動。", displayName)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			matches, err := getValorantMatches(name, tag)
			if err != nil {
				logger.Errorf("嘗試獲取玩家 %s 的遊戲狀態時出現錯誤: %v", displayName, err)
				continue
			} else if len(matches) == 0 {
				logger.Warnf("玩家 %s 的對戰記錄為空，已略過。", displayName)
			}

			latestData := matches[0].MetaData
			logger.Debugf("玩家 %s 的最新對戰資料ID為 %v, 時間: %s", displayName, latestData.MatchId, datetime.FormatSeconds(latestData.GameStart))

			lastMatchId, ok := lastMatchMap.Load(displayName)

			// 與上一次的狀態相同 => 忽略
			if ok && lastMatchId == latestData.MatchId {
				continue
			}

			lastMatchMap.Store(displayName, latestData.MatchId)

			// 尚未有上一次檢測的資料 + 最新對戰記錄距今已超過24小時 => 忽略
			if !ok && datetime.Duration(latestData.GameStart, time.Now().Unix()).Hours() > 24 {
				continue
			}

			publishData := &MatchMetaDataPublish{
				Data:        &latestData,
				DisplayName: fmt.Sprintf("%s#%s", name, tag),
			}

			publish(displayName, publishData)
		}
	}

}
