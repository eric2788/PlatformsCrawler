package twitter

import (
	"fmt"
	"strings"
	"github.com/eric2788/PlatformsCrawler/crawling"
)

const key = "twitter:user_caches"

func GetUserIdWithCache(screenNames []string) (map[string]string, error) {
	return crawling.LoadWithCache(key, GetUserId, IsNotExistUser, screenNames...)
}

func GetUserId(screenNames []string) (map[string]string, error) {

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
		profile, err := scraper.GetUserIDByScreenName(user)
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

		userMapAlt, err := GetUserIdWithCache(altScreenNames)

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
