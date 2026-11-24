package config

import "github.com/spf13/viper"

type AuthConfig struct {
	AppPort    string `mapstructure:"APP_PORT"`
	AppName    string `mapstructure:"APP_NAME"`
	SecretKey  string `mapstructure:"SECRET_KEY"`
	RefreshKey string `mapstructure:"REFRESH_KEY"`
	ExpireIn   int    `mapstructure:"EXPIRE_IN"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	GRPCPort   string `mapstructure:"GRPC_PORT"`

	//Redis
	RedisHost string `mapstructe:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`

	//Kafka
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
}

func LoadConfig() (*AuthConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config AuthConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil

}
