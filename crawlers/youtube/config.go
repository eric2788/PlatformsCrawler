package youtube

type Configuration struct {
	Interval        int64  `yaml:"interval"`
	NotLiveKeyword  string `yaml:"notLiveKeyword"`
	UpComingKeyword string `yaml:"upComingKeyword"`
	LiveKeyword     string `yaml:"liveKeyword"`

	Api *ApiConfig `yaml:"api"`
}

type ApiConfig struct {
	Key      string `yaml:"key"`
	Region   string `yaml:"region"`
	Language string `yaml:"language"`
}

var youtubeYaml = &Configuration{
	Interval:        60,
	NotLiveKeyword:  "{\"tabRenderer\":",
	UpComingKeyword: "\"isUpcoming\":true",
	LiveKeyword:     "\"videoViewCountRenderer\":",
	Api: &ApiConfig{
		Region:   "SG",
		Language: "zh-TW",
	},
}
