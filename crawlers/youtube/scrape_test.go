package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/eric2788/common-utils/regex"
	"github.com/eric2788/common-utils/request"
	"github.com/sirupsen/logrus"
)

var channels = map[string]string{
	"SasakiSaku": "UCoztvTULBYd3WmStqYeoHcA",
	"Chima":      "UCo7TRj3cS-f_1D9ZDmuTsjw",
	"Luna":       "UCpqNOJU70lNpa65FHkbrn6g",
	"Komori":     "UCBIR44irWpj1eTx0ZQFofHg",
	"Serena":     "UCRXBTd80F5IIWWY4HatJ5Ug",
	"music":      "UCcHWhgSsMBemnyLhg6GL1vA",
	"nano":       "UC0lIq8G4LgDPlXsDmYSUExw",
	"otto":       "UCvEX2UICvFAa_T6pqizC20g",
}

// TestGetChannelLiveResponse try to use with ytInitialData, referenced by Sora233/DDBOT
func TestGetChannelLiveResponse(t *testing.T) {

	if err := os.Mkdir("result", 0775); err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}

	for name, channel := range channels {

		func(name, channelId string) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://youtube.com/channel/%s/live", channelId), nil)

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")

			res, err := http.DefaultClient.Do(req)

			if err != nil {
				t.Fatal(err)
			} else if res.StatusCode == 404 {
				t.Fatal(err)
			}

			defer res.Body.Close()

			content, err := io.ReadAll(res.Body)

			if err != nil {
				t.Fatal(err)
			}

			if f, err := os.Create(fmt.Sprintf("result/%s_%s.html", name, channelId)); err == nil {
				f.Write(content)
				f.Close()
			} else {
				t.Fatal(err)
			}

			var reg *regexp.Regexp
			if strings.Contains(string(content), `window["ytInitialData"]`) {
				reg = regexp.MustCompile("window\\[\"ytInitialData\"\\] = (?P<json>.*);")
			} else {
				reg = regexp.MustCompile(">var ytInitialData = (?P<json>.*?);</script>")
			}

			result := regex.GetParams(reg, string(content))

			if jsonStr, ok := result["json"]; !ok {
				t.Log("cannot find ytInitialData")
			} else {

				var jsonObj = make(map[string]interface{})
				if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
					t.Log(err)
				}

				if b, err := json.MarshalIndent(jsonObj, "", "  "); err == nil {

					f, err := os.Create(fmt.Sprintf("result/%s_%s.json", name, channelId))

					if err != nil {
						t.Log(err)
					}

					f.Write(b)

					f.Close()

				} else {
					t.Log(err)
				}
			}

		}(name, channel)
	}
}

func TestGetOneChannelStatus(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	// load youtube yaml
	file.LoadYaml("youtube", youtubeYaml)
	initKeywordRegexp()

	name := "otto"
	id := channels[name]

	for i := 0; i < 5; i++ {
		status, err := getChannelStatusDebug(id, i)

		if err != nil {
			t.Fatalf("GetChannelStatus Error: %v", err)
		} else {
			if b, err := json.MarshalIndent(status, "", "\t"); err != nil {
				t.Fatalf("Json Marshal Error: %v", err)
			} else {
				fmt.Printf("%s 的直播狀態 \n", name)
				fmt.Println(string(b))

				if status.Id != "" && youtubeYaml.Api.Key != "" {
					if err = showVideoContent(status.Id); err != nil {
						t.Fatal(err)
					}
				}
			}
		}
	}
}

func TestGetChannelStatus(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)

	// load youtube yaml
	file.LoadYaml("youtube", youtubeYaml)

	initKeywordRegexp()

	for name, id := range channels {

		status, err := GetChannelStatus(id)

		if err != nil {
			t.Fatalf("GetChannelStatus Error: %v", err)
		} else {
			if b, err := json.MarshalIndent(status, "", "\t"); err != nil {
				t.Fatalf("Json Marshal Error: %v", err)
			} else {
				fmt.Printf("%s 的直播狀態 \n", name)
				fmt.Println(string(b))

				if status.Id != "" && youtubeYaml.Api.Key != "" {
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
