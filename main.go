package main

import (
	"flag"
	"monitor/common"
	"monitor/logger"
	"monitor/processor"
)

var (
	err      error
	confFile string
)

// 解析命令行参数
func initArgs() {
	// go run main.go -config config/config.yml -xxx 123 -yyy ddd
	flag.StringVar(&confFile, "config", "config.yml", "go run main.go -config ../config.yml")
	flag.Parse()
}

func main() {

	// 初始化命令行参数
	initArgs()

	if err = common.InitConfig(confFile); err != nil {
		goto ERR
	}

	if err = processor.InitReader(); err != nil {
		goto ERR
	}

	if err = processor.InitResolver(); err != nil {
		goto ERR
	}

	if err = processor.InitWriter(); err != nil {
		goto ERR
	}

	if err = processor.InitMonitor(); err != nil {
		goto ERR
	}

	go common.MockData()

	logger.Info.Printf("服务启动成功 ... ...")
	select {}

ERR:
	logger.Error.Printf("服务启动失败：%s", err)
}
