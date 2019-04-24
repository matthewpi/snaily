package backend

import (
	"github.com/go-redis/redis"
)

// RedisDriver represents a "stacktrace.fun" Redis driver.
type RedisDriver struct {
	Client *redis.Client
}

// Connect .
func (driver *RedisDriver) Connect(uri string, password string, database int) error {
	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: password,
		DB:       database,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	driver.Client = client

	return nil
}
