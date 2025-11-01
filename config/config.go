package config

import (
	"flag"
	"fmt"
	"time"
)

type MessageProcessorAPI struct {
	Host  string
	Token string
}

type Message struct {
	Retry             bool
	Timeout           time.Duration
	SendInterval      time.Duration
	BatchProcessCount int
	RetryFailCount    int
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
	Api              API
	Message          Message
	DB               DBConn
	Redis            RedisConn
	MessageProcessor MessageProcessorAPI
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
	interval := flag.Duration("interval", 2*time.Minute, "Interval between message sends") // put: 30s, 1m, 2m30s
	dbPassword := flag.String("password", "", "Password for db connection")
	redisPassword := flag.String("redis-password", "", "Password for redis db connection")
	batchProcessCount := flag.Int("message-process-count", 2, "Password for redis db connection")

	flag.Parse()

	cfg := Config{
		Message: Message{
			Retry:             *retry,
			Timeout:           *timeout,
			SendInterval:      *interval,
			BatchProcessCount: *batchProcessCount,
			RetryFailCount:    2, // keep it as 2 for now
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
			Port:     6379,
			Password: *redisPassword,
		},
		Api: API{
			Port: 8080,
		},
		MessageProcessor: MessageProcessorAPI{
			Host:  "https://webhook.site",
			Token: "75e83016-f99f-46ad-b833-8915a9abd327", // use this if you get 429 378d51ae-3661-4491-bc37-8da6434eb8c2, you need to chang
		},
	}

	return &cfg, nil
}
