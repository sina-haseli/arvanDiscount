package repositories

import (
	"github.com/go-redis/redis/v7"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *redisRepository {
	return &redisRepository{
		client: client,
	}
}

func (r *redisRepository) Increase(key string) (int, error) {
	res, err := r.client.Incr(key).Result()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *redisRepository) Decrease(key string) (int, error) {
	res, err := r.client.Decr(key).Result()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *redisRepository) SetValue(key string, value interface{}) error {
	_, err := r.client.Set(key, value, 0).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisRepository) GetValue(key string) (string, error) {
	v, err := r.client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return v, nil
}
