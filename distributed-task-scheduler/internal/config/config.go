package config

type Config struct {
	NumProducers     int
	NumWorkers       int
	TaskBufferSize   int
	ResultBufferSize int
}

func LoadConfig() Config {
	return Config{
		NumProducers:     3,
		NumWorkers:       5,
		TaskBufferSize:   100,
		ResultBufferSize: 100,
	}
}
