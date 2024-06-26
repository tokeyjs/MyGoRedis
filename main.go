package main

import (
	"MyGoRedis/config"
	"MyGoRedis/lib/logger"
	"MyGoRedis/resp/handler"
	"MyGoRedis/tcp"
	"encoding/json"
	"fmt"
	"os"
)

// redis配置文件
const configFile string = "redis.conf"

// 默认配置
var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func init() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "MyGoRedis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	// 从配置文件中载入配置
	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}
	logger.Infof("mygoredis config load successful!")
	cf, _ := json.Marshal(config.Properties)
	logger.Infof("config: %v", string(cf))
}

func main() {

	// 开启 服务
	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port)},
		handler.MakeHandler(),
	)
	if err != nil {
		logger.Errorf("tcp start error: %v", err)
	}

}
