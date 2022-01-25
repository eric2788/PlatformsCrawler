package bilibili

type Configuration struct {
	BiliLiveHost      string `yaml:"bili_live_host"`
	UseTLS            bool   `yaml:"use_tls"`
	AntiDuplicateLive int64  `yaml:"anti_duplicate_live"`
}

var bilibiliYaml = &Configuration{
	BiliLiveHost:      "blive.ericlamm.xyz",
	UseTLS:            true,
	AntiDuplicateLive: 10,
}
