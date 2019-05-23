package client

type Config struct {
	Host string
}

func DefaultConfig() Config {
	return Config{
		Host: "http://localhost:8080",
	}
}
