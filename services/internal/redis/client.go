package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var Client *redis.Client

func Init() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		opts, err := redis.ParseURL(redisURL)
		if err == nil {
			fmt.Println("Using Redis URL:", redisURL)
			Client = redis.NewClient(opts)
		} else {
			fmt.Println("Failed to parse REDIS_URL, falling back to manual connection")
		}
	}

	if Client == nil {
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			host = os.Getenv("REDISHOST")
		}
		if host == "" {
			host = "localhost"
		}

		port := os.Getenv("REDIS_PORT")
		if port == "" {
			port = os.Getenv("REDISPORT")
		}
		if port == "" {
			port = "6379"
		}

		password := os.Getenv("REDIS_PASSWORD")
		if password == "" {
			password = os.Getenv("REDISPASSWORD")
		}

		dbStr := os.Getenv("REDIS_DB")
		db := 0
		if dbStr != "" {
			if d, err := strconv.Atoi(dbStr); err == nil {
				db = d
			}
		}

		fmt.Printf("Connecting to Redis at %s:%s with password: %t\n", host, port, password != "")

		Client = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Password:     password,
			DB:           db,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     20,
			MinIdleConns: 5,
		})
	}

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		fmt.Println("Redis connection failed:", err)
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	fmt.Println("Connected to Redis successfully!")
	return nil
}

func SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

func GetCache(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}
