package tests

import (
	"assignment/config"
	"time"
)

func defaultConfig() *config.Config {
	return &config.Config{
		Message: config.Message{
			Retry:             true,
			Timeout:           5 * time.Second,
			SendInterval:      2 * time.Minute,
			BatchProcessCount: 2,
			RetryFailCount:    2,
		},
		WorkerPool: config.WorkerPool{
			Size:       3,
			BufferSize: 10,
		},
		DB: config.DBConn{
			Host:     "localhost",
			User:     "appuser",
			Port:     5432,
			Password: "secret",
			DBName:   "appdb",
		},
		Redis: config.RedisConn{
			Host:     "localhost",
			Port:     6379,
			Password: "",
		},
		Api: config.API{
			Port: 8080,
		},
		MessageProcessor: config.MessageProcessorAPI{
			Host:  "https://webhook.site",
			Token: "75e83016-f99f-46ad-b833-8915a9abd327", // use this if you get 429 378d51ae-3661-4491-bc37-8da6434eb8c2, you need to chang
		},
	}
}
