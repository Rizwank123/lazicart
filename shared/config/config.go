package config

import (
	"log"

	"github.com/spf13/viper"
)

type GlobalConfig struct {
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
}

func LoadConfig() (*GlobalConfig, error) {
	viper.AddConfigPath("../../shared") // Add path to shared folder
	viper.AddConfigPath(".")            // Also look in current directory
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Log but don't fail yet, maybe Env Vars are set in the system
		log.Println("Warning: .env file not found, checking environment variables")
	}

	var config GlobalConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Logic Fix: Assign default value BEFORE failing
	if len(config.KafkaBrokers) == 0 {
		log.Println("KAFKA_BROKERS not found in config, using default localhost:9092")
		config.KafkaBrokers = []string{"localhost:9092"}
	}

	return &config, nil

}
