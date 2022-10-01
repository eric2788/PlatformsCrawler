package valorant

import (
	"github.com/eric2788/common-utils/datetime"
	"testing"
)

func TestGetValorantMatches(t *testing.T) {
	data, err := getValorantMatches("suou", "9035")
	if err != nil {
		t.Log(err)
	} else {
		for _, d := range data {
			t.Log(datetime.FormatSeconds(d.MetaData.GameStart))
		}
	}
}
