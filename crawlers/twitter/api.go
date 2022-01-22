package twitter

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/eric2788/PlatformsCrawler/crawling"
	"strings"
)

var (
	key      = "twitter:user_caches"
	notFound = mapset.NewSet()
)

func UserLookUpCache(screenNames []string) (map[string]string, error) {

	var cache map[string]string

	logger.Debugf("prepare to lookup users from cache: %v", screenNames)

	if notInCaches, err := crawling.GetMap(key, &cache, screenNames...); err != nil {
		return nil, err
	} else if len(notInCaches) == 0 { // all names are in the cache
		return cache, nil
	} else { // some names can't get

		toFetch := make([]string, 0)

		// 透過快取作過濾
		for _, name := range notInCaches {
			if notFound.Contains(name) {
				continue
			}
			toFetch = append(toFetch, name)
		}

		if len(toFetch) == 0 {
			return cache, nil
		}

		fetched, err := UserLookUp(toFetch)

		if err != nil {
			// 錯誤返回 twitter error code 17 代表 所有 toFetch 結果均為無效
			if IsNotExistUser(err) {
				logger.Warnf("查無這些用戶: %v", strings.Join(toFetch, ", "))
				for _, screen := range toFetch {
					notFound.Add(screen)
				}
				return cache, nil
			}
			return nil, err
		}

		if fetched != nil {

			if err := crawling.SaveMap(key, fetched); err != nil {
				logger.Warnf("Redis 儲存 用戶資訊 時 出現錯誤: %v", err)
				logger.Warnf("儲存的資訊: %v", fetched)
			}

			for screen, id := range fetched {
				cache[screen] = id
			}

		}

		return cache, nil
	}
}

func UserLookUp(screenNames []string) (map[string]string, error) {

	if len(screenNames) == 0 {
		return nil, nil
	}

	var altScreenNames []string

	// max per request is 100 users
	if len(screenNames) > 100 {
		altScreenNames = screenNames[100:]
		screenNames = screenNames[:100]
	}

	logger.Debugf("prepare to lookup users: %v", screenNames)

	users, _, err := client.Users.Lookup(&twitter.UserLookupParams{
		ScreenName:      screenNames,
		IncludeEntities: twitter.Bool(false),
	})

	if err != nil {
		return nil, err
	}

	userMap := make(map[string]string)

	//把得出的結果加到 set 來檢查某些沒有加到結果的 screenNames
	nameSet := mapset.NewSet()

	for _, user := range users {
		userMap[user.ScreenName] = user.IDStr
		nameSet.Add(user.ScreenName)
	}

	// 然後加到 notFound 中作為例外
	for _, name := range screenNames {
		if !nameSet.Contains(name) {
			notFound.Add(name)
		}
	}

	if len(altScreenNames) > 0 {

		userMapAlt, err := UserLookUpCache(altScreenNames)

		if err != nil {
			return userMap, err
		}

		for screen, id := range userMapAlt {
			userMap[screen] = id
		}
	}

	return userMap, nil
}

func IsNotExistUser(err error) bool {
	twErr, ok := err.(twitter.APIError)
	return ok && twErr.Errors[0].Code == 17
}
