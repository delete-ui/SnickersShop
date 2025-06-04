package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env      string   `yaml:"env"`
	Server   Server   `yaml:"server"`
	Postgres Postgres `yaml:"postgres"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func LoadConfig(configName string) (*Config, error) {

	viper.SetConfigName(configName)
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg *Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil

}
