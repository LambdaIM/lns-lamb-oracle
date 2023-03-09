package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

const YamlFile = "config.yaml"

type Config struct {
	Market struct {
		Url    string `yaml:"url"`
		ApiKey string `yaml:"apiKey"`
	} `yaml:"market"`
	Lambda struct {
		RpcUrl     string `yaml:"rpcUrl"`
		DurationH  uint64 `yaml:"durationH"`
		PrivateKey string `yaml:"privateKey"`
		OracleAddr string `yaml:"oracleAddr"`
	} `yaml:"lambda"`
}

func NewConfig(pwd string) (*Config, error) {
	file := filepath.Join(pwd, "config", YamlFile)
	data, err := os.ReadFile(file)
	if nil != err {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if nil != err {
		return nil, err
	}
	return cfg, nil
}
