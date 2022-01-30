package file

import (
	"fmt"
	"github.com/eric2788/PlatformsCrawler/logging"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	logger                         = logging.GetMainLogger()
	ApplicationYaml *Configuration = &defaultAppConfig
)

type Configuration struct {
	// seconds
	CheckInterval int64               `yaml:"checkInterval"`
	Redis         *RedisConfiguration `yaml:"redis"`
}

type RedisConfiguration struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database int    `yaml:"database"`
	Password string `yaml:"password"`
}

var defaultAppConfig = Configuration{
	CheckInterval: 5,
	Redis: &RedisConfiguration{
		Host:     "192.168.0.127",
		Port:     6379,
		Database: 1,
		Password: "",
	},
}

func LoadApplicationYaml() {
	LoadYaml("application", ApplicationYaml)
}

func LoadYaml(name string, defaultData interface{}) {

	// 創建 config 文件夾
	if err := os.MkdirAll("config", 0775); err != nil {

		logger.Errorf("生成 config 資料夾時出現錯誤: %v", err)
		os.Exit(1)

	}

	yml := fmt.Sprintf("config/%s.yaml", name)

	content, readErr := os.ReadFile(yml)

	if readErr != nil {

		if os.IsNotExist(readErr) {

			logger.Infof("%s 不存在，正在生成默認的文件...", yml)

			if def, err := yaml.Marshal(defaultData); err != nil {

				logger.Errorf("獲取默認文件內容時出現錯誤: %v", err)
				os.Exit(1)

			} else {
				if err := os.WriteFile(yml, def, 0775); err != nil {
					logger.Errorf("生成默認文件時出現錯誤: %v", err)
					os.Exit(1)
				} else {
					logger.Infof("生成 %s 的默認文件成功", yml)
					content = def
				}
			}
		} else {
			logger.Errorf("讀取文件時出現未知錯誤: %v", readErr)
			os.Exit(1)
		}
	}

	if err := yaml.Unmarshal(content, defaultData); err != nil {
		logger.Errorf("讀取文件 %s 出現錯誤: %v", yml, err)
	}

}
