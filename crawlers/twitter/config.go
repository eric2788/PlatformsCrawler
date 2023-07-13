package twitter

type Configuration struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	EmailCode string `yaml:"email_code"`
	ScrapeInterval int `yaml:"scrape_interval"`
	RequestDelay int64 `yaml:"request_delay"`
}

var twitterYaml = &Configuration{
	ScrapeInterval: 60,
	RequestDelay: 5,
}
