package youtube

import (
	"encoding/json"
	"fmt"
	"github.com/eric2788/common-utils/request"
	"net/http"
	"regexp"
	"testing"
)

var channels = map[string]string{
	"SasakiSaku": "UCoztvTULBYd3WmStqYeoHcA",
	"Chima":      "UCo7TRj3cS-f_1D9ZDmuTsjw",
	"Luna":       "UCpqNOJU70lNpa65FHkbrn6g",
	"Komori":     "UCBIR44irWpj1eTx0ZQFofHg",
	"Serena":     "UCRXBTd80F5IIWWY4HatJ5Ug",
	"music":      "UCcHWhgSsMBemnyLhg6GL1vA",
	"nano":       "UC0lIq8G4LgDPlXsDmYSUExw",
}

func TestGetChannelStatus(t *testing.T) {

	// load youtube yaml
	instance.Init()

	initKeywordRegexp()

	for name, id := range channels {

		status, err := GetChannelStatus(id)

		if err != nil {
			t.Fatal(err)
		} else {
			if b, err := json.MarshalIndent(status, "", "\t"); err != nil {
				t.Fatal(err)
			} else {
				fmt.Printf("%s 的直播狀態 \n", name)
				fmt.Println(string(b))

				if status.Id != "" {
					if err = showVideoContent(status.Id); err != nil {
						t.Fatal(err)
					}
				}
			}
		}
	}
}

func TestDoubleReadFind(t *testing.T) {
	r := regexp.MustCompile("google")
	r2 := regexp.MustCompile("search")
	res, err := http.Get("https://google.com")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	c, err := request.ReadForRegex(res, r, r2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c[0])
	fmt.Println(c[1])
}

func showVideoContent(id string) error {

	if youtubeService == nil {
		initYoutubeService()
	}

	video, err := getVideos(id)

	if err != nil {
		return err
	}

	/*
		if b, err := json.MarshalIndent(video[0].Snippet, "", "\t"); err != nil {
			return err
		} else {
			fmt.Printf(string(b))
		}

	*/
	fmt.Println(video[0].Snippet.Title)
	return nil
}
