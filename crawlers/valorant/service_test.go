package valorant

import (
	"github.com/eric2788/common-utils/datetime"
	"testing"
)

func TestGetValorantMatches(t *testing.T) {
	data, err := getValorantMatches("f4a508ce-d7c3-561c-9d36-2d6808c18f10")
	if err != nil {
		t.Log(err)
	} else {
		for _, d := range data {
			t.Log(datetime.FormatSeconds(d.MetaData.GameStart))
		}
	}
}

func TestGetDisplayName(t *testing.T) {
	data, err := getDisplayName("f4a508ce-d7c3-561c-9d36-2d6808c18f10")
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("%s#%s", data.Name, data.Tag)
	}
}
