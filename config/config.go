package config

import (
	"github.com/spf13/viper"
)

const DefaultAddress = "127.0.0.1:8080"

type Config struct {
	ClickHouseHost                string `mapstructure:"CLICKHOUSE_HOST"`
	HTTPServerAddress             string `mapstructure:"HTTP_SERVER_ADDRESS"`
	VideoProcessingServiceAddress string `mapstructure:"VIDEO_PROCESSING_SERVICE_ADDRESS"`
	VideoIndexingServiceAddress   string `mapstructure:"VIDEO_INDEXING_SERVICE_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.SetDefault("HTTP_SERVER_ADDRESS", DefaultAddress)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
