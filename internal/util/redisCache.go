package util

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

func ConnectToRedisCache() (*redis.Client, func() error) {

	//rdsUri := os.Getenv("REDIS_CACHE_URI")
	//if rdsUri == "" {
	//	Logger.Fatal("Set your 'REDIS_CACHE_URI' environment variable. ")
	//}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// Ensure that you have Redis running on your system

	rdb := redis.NewClient(&redis.Options{
		Addr:       "localhost:6380",
		DB:         1,
		MaxRetries: 3,
	})

	// Ensure that the connection is properly closed gracefully

	// Perform basic diagnostic to check if the connection is working
	// Expected result > ping: PONG
	// If Redis is not running, error case is taken instead
	status, err := rdb.Ping(ctx).Result()

	if err != nil {
		Logger.Fatalln("Redis connection was refused")
	}
	Logger.Println(status)

	return rdb, rdb.Close
}
