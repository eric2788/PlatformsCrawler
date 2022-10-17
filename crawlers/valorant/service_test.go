package valorant

import (
	"sync"
	"testing"

	"github.com/eric2788/common-utils/datetime"
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

func TestGetValorantMatchesLoop(t *testing.T) {
	wg := &sync.WaitGroup{}
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup){
			defer w.Done()
			_, err := getValorantMatches("f4a508ce-d7c3-561c-9d36-2d6808c18f10")
			if err != nil {
				t.Log(err)
			} else {
				t.Log("no error")
			}
		}(wg)
	}

	wg.Wait()
}

func TestGetDisplayName(t *testing.T) {
	data, err := getDisplayName("f4a508ce-d7c3-561c-9d36-2d6808c18f10")
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("%s#%s", data.Name, data.Tag)
	}
}
