# Platform Crawlers

多平台爬蟲 + 模塊化管理，透過 Redis Pubsub 推送

## 目前內置爬蟲

### Youtube
- 預定直播
- 開播

### Twitter
- 推文/轉發推文
- 回復

### Bilibili
- 所有直播數據流 (with [biligo-live-ws](https://github.com/eric2788/biligo-live-ws))


## 原理

透過隔秒檢查 Redis pubsub 內符合 pattern 的 channels，經過去重過濾取出要新增監聽和中止監聽的頻道，然後基於機制:

- 為 `EachCrawling` 則各自啟動和中止
- 為 `OnceCrawling` 則先關閉再啟動以刷新

實現: [crawling.go](/crawling/crawling.go)


## 部署

你需要:
- redis 伺服器
- biligo-live-ws 以推送 bilibili 直播數據流

參數:
- `--debug` 啟動 debug level (默認: false)
- `--port` 啟動 rest api 的 port (默認: 8989)

Docker 部署: `docker.io/eric1008818/platform-crawlers`

## 使用

- 開啟 platform-crawlers 程序後關閉，設置好 `config` 內所有 yaml 再重開
- 首先把 platform-crawlers 與 監聽程序 連接到同一個 redis 伺服器
- 在 監聽程序 訂閱 格式為 `[prefix]:[room]` 的 topic (eg. `blive:22671795` 將監聽房間號為 22671795 的B站直播)
- 訂閱後，將會開始自動接收推送

## 新增新的推送

由於此程序採用模塊化管理，其新增新的推送極其簡單(需要使用 `golang`):

- 創建一個 `struct` 並實現 `Crawler`
- 根據你的監聽方式實現 `EachCrawling` 或 `OnceCrawling` 二選一
- 在 `init` 方法中使用 `crawling.RegisterCrawler(Tag, crawler instance, logger)`
- 最後，在 `main.go` 透過 `_` import 你的 package 即可 
- 如果想禁用某些推送，可以在 `application.yml` 中的 `disabledCrawlers` 屬性中加入你想要禁用的推送 Tag

### 此爬蟲目前主要負責用於我的私群專用機器人上，詳見 [mirai-val-bot](https://github.com/eric2788/miraivalbot)