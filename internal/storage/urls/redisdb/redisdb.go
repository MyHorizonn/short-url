package redisdb

import (
	"context"
	"short-url/internal/storage/urls"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func (r *Redis) Create(url urls.Url) error {
	ctx := context.Background()
	err := r.Client.Set(ctx, strconv.FormatUint(url.Key, 10), url.OriginalUrl, time.Duration(time.Hour*24*7)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(key uint64) (string, error) {
	ctx := context.Background()
	val, err := r.Client.Get(ctx, strconv.FormatUint(key, 10)).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Redis) IsExists(key uint64) (bool, error) {
	ctx := context.Background()
	_, err := r.Client.Get(ctx, strconv.FormatUint(key, 10)).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
