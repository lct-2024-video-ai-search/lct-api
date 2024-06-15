package config

import (
	"github.com/spf13/viper"
	"log"
)

const DefaultAddress = "127.0.0.1:8080"

type Config struct {
	Environment                   string `mapstructure:"ENVIRONMENT"`
	DBDriver                      string `mapstructure:"DB_DRIVER"`
	DBSource                      string `mapstructure:"DB_SOURCE"`
	ClickHouseHost                string `mapstructure:"CLICKHOUSE_HOST"`
	MigrationURL                  string `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress             string `mapstructure:"HTTP_SERVER_ADDRESS"`
	VideoProcessingServiceAddress string `mapstructure:"VIDEO_PROCESSING_SERVICE_ADDRESS"`
	VideoIndexingServiceAddress   string `mapstructure:"VIDEO_INDEXING_SERVICE_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	if config.HTTPServerAddress == "" {
		log.Printf("Server address is unset, assuming %s", DefaultAddress)
		config.HTTPServerAddress = DefaultAddress
	}
	return
}
