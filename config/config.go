package config

import "github.com/spf13/viper"

type Config struct {
	EnableGroup bool    `yaml:"enable-group"`
	Masters     []int64 `yaml:"masters"`
}

func LoadConfig() (*Config, error) {
	var config Config
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./JarvisConfig.yml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
