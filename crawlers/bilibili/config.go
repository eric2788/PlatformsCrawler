package bilibili

type Configuration struct {
	BiliLiveHost string `yaml:"bili_live_host"`
	UseTLS       bool   `yaml:"use_tls"`
}

var bilibiliYaml = &Configuration{
	BiliLiveHost: "blive.ericlamm.xyz",
	UseTLS:       true,
}
