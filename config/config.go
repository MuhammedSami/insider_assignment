package config

import (
	"flag"
	"fmt"
	"time"
)

type Message struct {
	Retry        bool
	Timeout      time.Duration
	SendInterval time.Duration
}

type DBConn struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

type RedisConn struct {
	Host     string
	Port     int
	Password string
}

type API struct {
	Port int
}
type Config struct {
	Api     API
	Message Message
	DB      DBConn
	Redis   RedisConn
}

// use this if config needs any other validation
func (c *Config) Validate() error {
	if c.DB.Password == "" {
		return fmt.Errorf("DB password is required")
	}

	return nil
}

func NewConfig() (*Config, error) {
	// we could use yaml based default approach and secret manager to load sensible information, but I will use flags for now

	retry := flag.Bool("retry", false, "Enable retry for failed messages")
	timeout := flag.Duration("timeout", 5*time.Second, "Timeout for each message send attempt")
	interval := flag.Duration("interval", 2*time.Minute, "Interval between message sends")
	dbPassword := flag.String("password", "", "Password for db connection")
	redisPassword := flag.String("redis-password", "", "Password for redis db connection")

	flag.Parse()

	cfg := Config{
		Message: Message{
			Retry:        *retry,
			Timeout:      *timeout,
			SendInterval: *interval,
		},
		DB: DBConn{
			Host:     "localhost",
			User:     "appuser",
			Port:     5432,
			Password: *dbPassword,
			DBName:   "appdb",
		},
		Redis: RedisConn{
			Host:     "localhost",
			Port:     6397,
			Password: *redisPassword,
		},
		Api: API{
			Port: 8080,
		},
	}

	return &cfg, nil
}
