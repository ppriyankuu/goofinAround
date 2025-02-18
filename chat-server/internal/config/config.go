package config

type Config struct {
	DBURL         string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func NewConfig() *Config {
	return &Config{
		DBURL:         "mongodb://localhost:27017",
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
	}
}
