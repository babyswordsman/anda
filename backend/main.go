package main

import (
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	logger.SetReportCaller(true)

	var configPath string
	flag.StringVar(&configPath, "conf", "conf/config.yaml", "config path,example:conf/config.yaml")
	flag.Parse()

	content, err := os.ReadFile(configPath)
	if err != nil {
		logger.Errorf("read file(%s) err:%s", configPath, err)
		return
	}
	config := &Config{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		logger.Errorf("parse yaml err:%s", err)
		return
	}

	if config.LogLevel != "" {
		lvl, err := logger.ParseLevel(config.LogLevel)
		if err != nil {
			logger.Errorf("parse log level err:%s", err.Error())
		}
		logger.SetLevel(lvl)
	}

	//star http server
	r := gin.Default()
	router(r)
	r.Run(config.ServerAddr)
	logger.Infof("star server:%s", config.ServerAddr)
}
