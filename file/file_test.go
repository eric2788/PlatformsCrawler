package file

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestYaml struct {
	A string `yaml:"a"`
	B int64  `yaml:"b"`
	C *struct {
		D string `yaml:"d"`
		E int64  `yaml:"e"`
		F *struct {
			G string `yaml:"g"`
		} `yaml:"f"`
	} `yaml:"c"`
}

var defaultTestData = TestYaml{
	A: "AAAA",
	B: 1234,
	C: &struct {
		D string `yaml:"d"`
		E int64  `yaml:"e"`
		F *struct {
			G string `yaml:"g"`
		} `yaml:"f"`
	}{
		D: "DDDD",
		E: 6666,
		F: &struct {
			G string `yaml:"g"`
		}{
			G: "GGGGG",
		},
	},
}

func TestLoadYaml(t *testing.T) {
	LoadYaml("test", &defaultTestData)
	assert.Equal(t, "GGGGG", defaultTestData.C.F.G)
}

func TestSet(t *testing.T) {
	a := mapset.NewSet(1, 2, 3)
	b := mapset.NewSet("a", "b", "c")

	fmt.Println(a.String())
	a = b
	fmt.Println(a.String())

}
