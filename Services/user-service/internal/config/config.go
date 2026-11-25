package config

import "github.com/spf13/viper"

type UserConfig struct {
	APPName      string   `mapstructure:"app_name"`
	AppPort      string   `mapstructure:"app_port"`
	DBHost       string   `mapstructure:"db_host"`
	DBPort       string   `mapstructure:"db_port"`
	DBUser       string   `mapstructure:"db_user"`
	DBPassword   string   `mapstructure:"db_password"`
	DBName       string   `mapstructure:"db_name"`
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
}

func LoadConfig() (*UserConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config UserConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil

}
