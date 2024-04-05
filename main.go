package main

import (
	"MyGoRedis/config"
	"MyGoRedis/resp/handler"
	"MyGoRedis/tcp"
	"fmt"
	"github.com/sirupsen/logrus"
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
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		PrettyPrint:       true,
		DisableHTMLEscape: false,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}
}

func main() {

	// 开启
	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port)},
		handler.MakeHandler(),
	)
	if err != nil {
		logrus.Errorf("tcp start error: %v", err)
	}

}
