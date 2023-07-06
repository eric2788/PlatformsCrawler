package twitter

type Configuration struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ScrapeInterval int `yaml:"scrape_interval"`
}

var twitterYaml = &Configuration{
	ScrapeInterval: 60,
}
