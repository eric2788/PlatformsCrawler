package twitter

type Configuration struct {
	BearerToken       string `yaml:"bearer_token"`
	ConsumerKey       string `yaml:"consumer_key"`
	ConsumerSecret    string `yaml:"consumer_secret"`
	AccessToken       string `yaml:"access_token"`
	AccessTokenSecret string `yaml:"access_token_secret"`
}

var twitterYaml = &Configuration{}
