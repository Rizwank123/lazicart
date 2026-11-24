package config

import "github.com/spf13/viper"

type PaymentConfig struct {
	APPName    string `mapstructure:"APP_NAME"`
	AppPort    string `mapstructure:"APP_PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig() (*PaymentConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg PaymentConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "7706"
	}
	return &cfg, nil
}
