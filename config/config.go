package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	EnableGroup      bool    `mapstructure:"enable-group"`
	Masters          []int64 `mapstructure:"masters"`
	WorkingDirectory string  `mapstructure:"working-directory"`
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

	// manually check some fields of config
	// WorkingDirectory (wd)
	// wd not specified
	if len(config.WorkingDirectory) == 0 {
		config.WorkingDirectory, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	// back slash of wd
	if config.WorkingDirectory[len(config.WorkingDirectory)-1] != '/' {
		config.WorkingDirectory += "/"
	}

	return &config, nil
}

func (cfg *Config) PrintConfig() {
	fmt.Println("EnableGroup")
	fmt.Println("\t", cfg.EnableGroup)
	fmt.Println("Masters")
	for _, v := range cfg.Masters {
		fmt.Println("\t", v)
	}
	fmt.Println("WorkingDirectory")
	fmt.Println("\t", cfg.WorkingDirectory)
}
