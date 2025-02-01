package config

import "github.com/spf13/viper"

type Config struct {
	Kafka    KafkaConfig
	Redis    RedisConfig
	Postgres PostgresConfig
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type PostgresConfig struct {
	DSN string
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
