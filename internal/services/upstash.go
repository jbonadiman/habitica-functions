package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const counterKey = "habitica:pomodoro_count"

type UpstashDB struct {
	*redis.Client
}

func NewRedisClient(ctx context.Context, url string) *UpstashDB {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	opt.ConnMaxIdleTime = 5 * time.Minute
	client := redis.NewClient(opt)

	err = client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	return &UpstashDB{client}
}

func (r *UpstashDB) TickCounter(ctx context.Context) (int, error) {
	count, err := r.Incr(ctx, counterKey).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		year, month, day := time.Now().UTC().Add(24 * time.Hour).Date()
		resp, err := r.ExpireAt(
			ctx, counterKey, time.Date(
				day,
				month,
				year,
				0,
				0,
				0,
				0,
				time.UTC,
			),
		).Result()
		if err != nil {
			return 0, err
		}

		log.Println("expiration set:", resp)
		if !resp {
			return 0, errors.New("failed to set expiration")
		}
	}

	return int(count), nil
}

func (r *UpstashDB) ResetCounter(ctx context.Context) error {
	_, err := r.Del(ctx, counterKey).Result()
	if err != nil {
		return err
	}

	return nil
}
