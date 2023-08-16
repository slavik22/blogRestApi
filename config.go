package blogRestApi

import (
	"github.com/spf13/viper"
)

// Config is a config :)
type Config struct {
	HTTPAddr string `mapstructure:"HTTP_ADDRESS"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
	DBSource string `mapstructure:"DB_SOURCE"`
}

var (
	config Config
)

// Get reads config from environment. Once.
func Get(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	return &config, nil
}
