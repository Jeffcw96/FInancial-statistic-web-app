package db

import (
	"time"

	"github.com/go-redis/redis"
)

var Client = redis.NewClient(&redis.Options{
	//Address 127.0.0.1:6379
	Addr:        "redis-10908.c54.ap-northeast-1-2.ec2.cloud.redislabs.com:10908",
	Password:    "LunjeffdevsLife961117~Tai",
	DB:          0,
	ReadTimeout: 10 * time.Minute,
})
