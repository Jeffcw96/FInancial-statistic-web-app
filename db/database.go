package db

import (
	"time"

	"github.com/go-redis/redis"
)

var Client = redis.NewClient(&redis.Options{
	//Address 127.0.0.1:6379
	Addr:        "127.0.0.1:6379",
	Password:    "",
	DB:          0,
	ReadTimeout: 10 * time.Minute,
})
