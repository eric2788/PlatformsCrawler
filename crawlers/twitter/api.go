package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/eric2788/PlatformsCrawler/crawling"
)

var key = "twitter:user_caches"

func UserLookUpCache(screenNames []string) (map[string]string, error) {

	var cache map[string]string

	logger.Debugf("prepare to lookup users from cache: %v", screenNames)

	if er, err := crawling.GetMap(key, &cache, screenNames...); err != nil {
		return nil, err
	} else if len(er) == 0 { // all names are in the cache
		return cache, nil
	} else { // some names can't get

		fetched, err := UserLookUp(er)

		if err != nil {
			return cache, err
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
		logger.Error(err)
		return nil, err
	}

	userMap := make(map[string]string)

	for _, user := range users {
		userMap[user.ScreenName] = user.IDStr
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
