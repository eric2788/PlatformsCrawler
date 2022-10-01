package valorant

type Configuration struct {
	HenrikApiKey       string `yaml:"henrik_api_key"`
	Interval           int    `yaml:"interval"`              // crawling interval
	MaxAccountPerTimes int    `yaml:"max_account_per_times"` // max account to crawl per each
	Region             string `yaml:"region"`                // the region to crawl
}

var valorantYaml = &Configuration{
	Interval:           60,
	MaxAccountPerTimes: 50,
	Region:             "ap",
}
