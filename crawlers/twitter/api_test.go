package twitter

import (
	"github.com/eric2788/PlatformsCrawler/crawling"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/eric2788/PlatformsCrawler/rest"
	"github.com/eric2788/common-utils/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

type twitterResp struct {
	Exist bool              `json:"exist"`
	Data  map[string]string `json:"data"`
}

func TestRestful(t *testing.T) {

	file.LoadYaml("twitter", twitterYaml)
	crawling.InitRedis()

	go rest.StartServe(8989)

	resp := twitterResp{}

	if err := request.Get("http://127.0.0.1:8989/twitter/userExist/mi_tagun", &resp); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, resp.Exist)
	assert.Equal(t, "1493253948", resp.Data["id"])
	assert.Equal(t, "mi_tagun", resp.Data["screen_name"])
}

func TestGetNickName(t *testing.T) {
	// crawling.InitRedis()

	err := scraper.LoginOpenAccount()

	if err != nil {
		t.Skip(err)
	}

	p, err := scraper.GetProfile("mi_tagun")

	if err != nil {
		t.Skip(err)
	}

	t.Logf("%+v", p)
}
