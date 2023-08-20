package redisq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/wei-zero/wz/broker"
)

func New(rc *redis.Client, env string) (broker.MessageBroker, error) {
	return &redisPubSub{rc: rc, env: env}, nil
}

type redisPubSub struct {
	rc  *redis.Client
	env string
}

func (r redisPubSub) Publish(topic string, msg interface{}) error {
	v, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.rc.Publish(context.Background(), fmt.Sprintf("%s:%s", r.env, topic), v).Err()
}

func (r redisPubSub) Subscribe(topic string, hf broker.HandleFunc) (int, error) {
	ch := r.rc.Subscribe(context.TODO(), fmt.Sprintf("%s:%s", r.env, topic)).Channel()
	go func() {
		for message := range ch {
			err := hf(topic, []byte(message.Payload))
			if err != nil {
				slog.Error("redisPubSub: handle event error", "topic", topic, "payload", message.Payload, "err", err)
			}
		}
		slog.Info("redisPubSub: subscribe goroutine end")
	}()

	return 0, nil
}

func (r redisPubSub) Unsubscribe(topic string, id int) error {
	return nil
}

func (r redisPubSub) Close() error {
	slog.Info("redisPubSub: close")
	return r.rc.Close()
}
