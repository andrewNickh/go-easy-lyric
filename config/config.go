package config

import (
	"github.com/spf13/viper"
	"log"
)

var Instance Config

type Config struct {
	Host              string      `yaml:"Host"`
	Username          string      `yaml:"Username"`
	Password          string      `yaml:"Password"`
	MySqlUrl          string      `yaml:"MySqlUrl"`
	MySqlMaxIdle      int         `yaml:"MySqlMaxIdle"`
	MySqlMaxOpen      int         `yaml:"MySqlMaxOpen"`
	SlaveMySqlUrl     string      `yaml:"SlaveMySqlUrl"`
	SlaveMySqlMaxIdle int         `yaml:"SlaveMySqlMaxIdle"`
	SlaveMySqlMaxOpen int         `yaml:"SlaveMySqlMaxOpen"`
	BaseURL           string      `yaml:"BaseURL"`
	SearchBaseURL     string      `yaml:"SearchBaseURL"`
	RedisCache        RedisConfig `yaml:"RedisCache"`
}

type RedisConfig struct {
	Host      []string `yaml:"Host"`      // 连接地址
	Password  string   `yaml:"Password"`  // 密码
	DB        int      `yaml:"DB"`        // 库索引
	MaxIdle   int      `yaml:"MaxIdle"`   // 最大空闲数
	MaxActive int      `yaml:"MaxActive"` // 最大连接数
}

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err = viper.Unmarshal(&Instance)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
