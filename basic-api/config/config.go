package config

import (
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func AppPort() string {
	return viper.GetString("app.port")
}

func JWTSecret() string {
	return viper.GetString("auth.jwt_secret")
}

func RateLimit() int {
	return viper.GetInt("rate_limit")
}

func CORSAllowedOrigin() string {
	return viper.GetString("cors.allowed_origin")
}
