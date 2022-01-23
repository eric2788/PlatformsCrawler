package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/eric2788/PlatformsCrawler/crawling"
)

const key = "twitter:user_caches"

func UserLookUpCache(screenNames []string) (map[string]string, error) {
	return crawling.LoadWithCache(key, UserLookUp, IsNotExistUser, screenNames...)
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

// IsNotExistUser 錯誤返回 twitter error code 17 代表 所有 toFetch 結果均為無效
func IsNotExistUser(err error) bool {
	twErr, ok := err.(twitter.APIError)
	return ok && twErr.Errors[0].Code == 17
}
