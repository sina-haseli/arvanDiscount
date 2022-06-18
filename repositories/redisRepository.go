package repositories

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *redisRepository {
	return &redisRepository{
		client: client,
	}
}

func (r *redisRepository) Dequeue(channelName string) (string, error) {
	res, err := r.client.BLPop(0*time.Second, channelName).Result()
	if err != nil {
		return "", err
	}

	return res[0], nil
}

func (r *redisRepository) Enqueue(message []byte, channelName string) error {
	_, err := r.client.RPush(channelName, string(message)).Result()
	return err
}
