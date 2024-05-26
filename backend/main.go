package main

import (
	"flag"
	"github.com/anda-ai/anda/conf"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	logger.SetReportCaller(true)

	var configPath string
	flag.StringVar(&configPath, "conf", "conf/config.yml", "config path,example:conf/config.yml")
	flag.Parse()

	content, err := os.ReadFile(configPath)
	if err != nil {
		path, e := filepath.Abs(configPath)
		if e != nil {
			logger.Errorf("read file(%s) err:%s stat_err:%s", configPath, err, e)
			return
		}
		logger.Errorf("read file(%s) err:%s", path, err)
		return
	}
	config := &conf.Config{}
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

	if err := r.Run(config.ServerAddr); err != nil {
		logger.Errorf("star server err:%s", err)
		return
	}q
	logger.Infof("star server:%s", config.ServerAddr)
}
