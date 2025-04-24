package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Configs struct {
	DatabaseConfig struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
		TimeZone string `yaml:"timezone"`
	} `yaml:"database"`
	Server struct {
		Host         string        `yaml:"host"`
		Port         int           `yaml:"port"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
	} `yaml:"server"`
}

func LoadConfig(path string) (*Configs, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Configs
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
