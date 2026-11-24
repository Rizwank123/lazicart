package config

import "github.com/spf13/viper"

type CatalogConfig struct {
	APPName    string `mapstructure:"APP_NAME"`
	AppPort    string `mapstructure:"APP_PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig() (*CatalogConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config CatalogConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil

}
