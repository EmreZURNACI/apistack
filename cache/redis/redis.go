package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Handler struct {
	client *redis.Client
}

type Message struct {
	Key      []byte        `json:"key"`
	Value    []byte        `json:"value"`
	Duration time.Duration `json:"duration"`
}

func Connection() (*Handler, error) {

	var dsn string = fmt.Sprintf("redis://%s:%s@%s:%d/%d",
		viper.GetString("redis.user"),
		viper.GetString("redis.password"),
		viper.GetString("redis.hostname"),
		viper.GetInt("redis.port"),
		viper.GetInt("redis.db"))

	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, ErrURLParseFailed
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, ErrConnectionFailed
	}

	return &Handler{
		client: client,
	}, nil
}

func (h *Handler) Set(ctx context.Context, msg Message) error {
	//Duration=0 forever
	err := h.client.Set(ctx, string(msg.Key), string(msg.Value), msg.Duration).Err()
	if err != nil {
		return ErrSetDataFailed
	}
	return nil
}

func (h *Handler) Get(ctx context.Context, key string) ([]byte, error) {

	value, err := h.client.Get(ctx, key).Result()
	if err != nil {
		return nil, ErrGetDataFailed
	}
	return []byte(value), nil
}
