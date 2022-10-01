package valorant

import "testing"

func TestGetValorantMatches(t *testing.T) {
	data, err := getValorantMatches("suou", "9035")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(data)
}
