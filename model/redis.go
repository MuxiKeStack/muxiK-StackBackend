package model

import (
	"github.com/go-redis/redis"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type Redis struct {
	Self *redis.Client
}

type PubSub struct {
	Self *redis.PubSub
}

var (
	RedisDB        *Redis
	GradeSubClient *PubSub
)

const (
	GradeChan = "grade"
	MsgChan   = "msg"
)

func (*Redis) Init() {
	RedisDB = &Redis{
		Self: OpenRedisClient(),
	}
}

func (*Redis) Close() {
	_ = RedisDB.Self.Close()
}

func (*PubSub) Init(channel string) {
	GradeSubClient = &PubSub{
		Self: OpenRedisPubSubClient(channel),
	}
}

func (*PubSub) Close() {
	_ = GradeSubClient.Self.Close()
}

func OpenRedisClient() *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       0,
	})

	if _, err := r.Ping().Result(); err != nil {
		log.Fatal("Open redis failed", err)
	}
	return r
}

func OpenRedisPubSubClient(channel string) *redis.PubSub {
	return OpenRedisClient().Subscribe(channel)
}

func PublishMsg(msg []byte, ch string) error {
	return RedisDB.Self.Publish(ch, msg).Err()
}
