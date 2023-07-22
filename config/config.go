package config

import "github.com/spf13/viper"

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
	// back slash of WorkingDirectory
	if len(config.WorkingDirectory) == 0 {
		config.WorkingDirectory += "./"
	} else if config.WorkingDirectory[len(config.WorkingDirectory)-1] != '/' {
		config.WorkingDirectory += "/"
	}

	return &config, nil
}
