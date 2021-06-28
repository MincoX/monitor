package common

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"monitor/logger"
)

var (
	GConfig *Config
)

type Config struct {
	Influx    Influx
	Http      Http
	Processor Processor
}

type Influx struct {
	Addr        string `yaml:"addr"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Database    string `yaml:"database"`
	Precision   string `yaml:"precision"`
	Measurement string `yaml:"measurement"`
}

type Processor struct {
	LogPath    string `yaml:"log_path"`
	ReaderNum   int    `yaml:"reader_num"`
	ResolverNum int    `yaml:"resolver_num"`
	WriterNum   int    `yaml:"writer_num"`
}

type Http struct {
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

func InitConfig(filename string) (err error) {
	var (
		conf    Config
		content []byte
	)
	if content, err = ioutil.ReadFile(filename); err != nil {
		logger.Error.Printf("读取配置问文件失败：%s", err)
		return
	}

	// 2, 做JSON反序列化
	if err = yaml.Unmarshal(content, &conf); err != nil {
		logger.Error.Printf("反序列化配置文件失败：%s", err)
		return
	}
	GConfig = &conf
	return
}
