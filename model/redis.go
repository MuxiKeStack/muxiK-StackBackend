package model

import (
	"time"

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
		log.Error("Open redis failed. Connect failed.", err)
	}
	return r
}

func OpenRedisPubSubClient(channel string) *redis.PubSub {
	return OpenRedisClient().Subscribe(channel)
}

// ----------------------------------------------------------------------------

func PublishMsg(msg []byte, ch string) error {
	return RedisDB.Self.Publish(ch, msg).Err()
}

func SetKey(key string, value interface{}, expiration time.Duration) error {
	return RedisDB.Self.Set(key, value, expiration).Err()
}

func GetKey(key string) (string, error) {
	return RedisDB.Self.Get(key).Result()
}

func HashSet(key, field string, value interface{}) error {
	return RedisDB.Self.HSet(key, field, value).Err()
}

func HashGet(key, field string) (string, bool, error) {
	val, err := RedisDB.Self.HGet(key, field).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	return val, true, err
}

func HashExist(key, field string) (bool, error) {
	return RedisDB.Self.HExists(key, field).Result()
}
