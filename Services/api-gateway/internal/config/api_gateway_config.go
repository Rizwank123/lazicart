package config

import "github.com/spf13/viper"

type ApiGateWayConfig struct {
	AppPort                  string `mapstructure:"APP_PORT"`
	AuthServiceURL           string `mapstructure:"AUTH_SERVICE_URL"`
	UserServiceURL           string `mapstructure:"USER_SERVICE_URL"`
	ProductServiceURL        string `mapstructure:"PRODUCT_SERVICE_URL"`
	InventoryServiceURL      string `mapstructure:"INVENTORY_SERVICE_URL"`
	OrderServiceURL          string `mapstructure:"ORDER_SERVICE_URL"`
	CartServiceURL           string `mapstructure:"CART_SERVICE_URL"`
	PaymentServiceURL        string `mapstructure:"PAYMENT_SERVICE_URL"`
	NotificationServiceURL   string `mapstructure:"NOTIFICATION_SERVICE_URL"`
	DeliveryServiceURL       string `mapstructure:"DELIVERY_SERVICE_URL"`
	ReviewServiceURL         string `mapstructure:"REVIEW_SERVICE_URL"`
	WishlistServiceURL       string `mapstructure:"WISHLIST_SERVICE_URL"`
	SearchServiceURL         string `mapstructure:"SEARCH_SERVICE_URL"`
	RecommendationServiceURL string `mapstructure:"RECOMMENDATION_SERVICE_URL"`
	AnalyticsServiceURL      string `mapstructure:"ANALYTICS_SERVICE_URL"`
	AdminServiceURL          string `mapstructure:"ADMIN_SERVICE_URL"`
	JWTSecret                string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (*ApiGateWayConfig, error) {
	viper.SetConfigFile(".env") // <-- direct file
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg ApiGateWayConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
