package twitter

import (
	"fmt"
	"strings"
	"github.com/eric2788/PlatformsCrawler/crawling"
)

const (
	userIdKey = "twitter:user_caches"
    nickNameKey = "twitter:username_caches"
)

func GetUserIdWithCache(screenNames []string) (map[string]string, error) {
	return crawling.LoadWithCache(userIdKey, GetUserId, IsNotExistUser, screenNames...)
}

func GetUserId(screenNames []string) (map[string]string, error) {
	return getStringMultiple(
		screenNames,
		scraper.GetUserIDByScreenName,
		GetUserIdWithCache,
	)
}

func GetNickNameWithCache(screenNames []string) (map[string]string, error) {
	return crawling.LoadWithCache(nickNameKey, GetNickNames, IsNotExistUser, screenNames...)
}

func GetNickNames(screenNames []string) (map[string]string, error) {
	return getStringMultiple(
		screenNames,
		func(s string) (string, error) {
			profile, err := scraper.GetProfile(s)
			if err != nil {
				return "", err
			}
			return profile.Name, nil
		},
		GetNickNameWithCache,
	)
}

func getStringMultiple(screenNames []string, getter func(string)(string, error), cacheGetter func([]string)(map[string]string, error)) (map[string]string, error) {
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

	users := make(map[string]string)
	errors := make(map[string]error)

	for _, user := range screenNames {
		profile, err := getter(user)
		if err != nil {
			logger.Errorf("error while fetching user id for username %s: %v", user, err)
			errors[user] = err
		} else {
			users[user] = profile
		}
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("%v", errors)
	}

	userMap := make(map[string]string)

	if len(altScreenNames) > 0 {

		userMapAlt, err := cacheGetter(altScreenNames)

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
	return strings.Contains(err.Error(), "does not exist")
}
