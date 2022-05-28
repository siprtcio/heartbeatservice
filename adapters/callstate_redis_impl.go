package adapters

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)

type CallState struct {
	client *redis.Client
}

func (cs *CallState) InitCallState() error {

	redisAddr := fmt.Sprintf("%s", beego.AppConfig.String("redis_cache"))
	db, err := beego.AppConfig.Int("redis_db")

	if err != nil {
		return err
	}

	cs.client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: beego.AppConfig.String("redis_password"), // no password set
		DB:       db,                                       // use default DB
	})

	if cs.client == nil {
		return errors.New("cs.client is nil")
	}
	return nil
}

func (cs *CallState) InitCallStateInternal(redisAddr string, DB int, redisPass string) error {
	cs.client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass, // no password set
		DB:       DB,        // use default DB
	})

	if cs.client == nil {
		return errors.New("cs.client is nil")
	}
	return nil
}

func (cs CallState) GetCallState(callUUID string) ([]byte, error) {
	val, err := cs.client.Get(callUUID).Bytes()
	if err != nil {
		return []byte("UNKNOWN"), err
	}
	return val, nil
}

func (cs CallState) SetCallState(call_uuid string, state []byte) error {
	err := cs.client.Set(call_uuid, state, 0).Err()
	return err
}

func (cs CallState) DelCallState(key string) error {
	err := cs.client.Del(key).Err()
	return err
}

func (cs CallState) SetCallRoute(call_uuid string, route string) error {
	err := cs.client.Set("Route"+call_uuid, route, 0).Err()
	return err
}

func (cs CallState) SetCallRouteTimeout(call_uuid string, route string) error {
	err := cs.client.Set("Route"+call_uuid, route, time.Hour*12).Err()
	return err
}

func (cs CallState) GetCallRoute(callUUID string) string {
	val := cs.client.Get("Route" + callUUID).Val()
	return val
}

func (cs CallState) DeleteCallRoute(callUUID string) {
	cs.client.Del("Route" + callUUID)
}

func (cs CallState) SetHeartBeatCounter(callUUID string) {
	cs.client.Set("HeartbeatCount"+callUUID, 0, 0)
}

func (cs CallState) IncHeartBeatCounter(callUUID string) int64 {
	result, err := cs.client.Incr("HeartbeatCount" + callUUID).Result()
	if err != nil {
		beego.Error("IncHeartBeatCounter - Failed", err.Error())
	}
	return result
}

func (cs CallState) DeleteHeartBeatCounter(callUUID string) {
	cs.client.Del("HeartbeatCount" + callUUID)
}

func (cs CallState) SetCallQueue(callQueue string) error {
	_, err := cs.client.SAdd("VOICE_CALL_QUEUE", callQueue).Result()
	return err
}

func (cs CallState) GetCallQueue() []string {
	return cs.client.SMembers("VOICE_CALL_QUEUE").Val()
}
