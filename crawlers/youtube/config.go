package youtube

type Configuration struct {
	Interval        int64  `yaml:"interval"`
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
	UpComingKeyword: `"isUpcoming":true`,
	LiveKeyword:     `<link rel="canonical" href="https://www.youtube.com/watch\?v=(?P<id>\w+)">`,
	Api: &ApiConfig{
		Region:   "SG",
		Language: "zh-TW",
	},
}
