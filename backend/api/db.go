package api

import (
	"context"
	"github.com/go-redis/redis/v9"
)

type RedisDB struct {
	client *redis.Client
	Ctx    context.Context
}

func (db *RedisDB) Connect(addr string) error {
	db.client = redis.NewClient(&redis.Options{
		Addr:     addr, // use default Addr
		Password: "",   // no password set
		DB:       0,    // use default DB
	})
	err := db.client.Ping(db.Ctx).Err()
	return err
}

func (db *RedisDB) Exists(continent string) (bool, error) {
	res, err := db.client.Exists(db.Ctx, continent).Result()

	if err != nil {
		return false, err
	} else if res == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (db *RedisDB) Add(continent string, score int, category string) error {
	_, err := db.client.ZAdd(
		db.Ctx, continent, redis.Z{Score: float64(score), Member: category}).Result()
	return err
}

func (db *RedisDB) GetTop(continent string) ([]string, error) {
	result, err := db.client.ZRevRangeByScore(
		db.Ctx, continent, &redis.ZRangeBy{Max: "+inf", Min: "0"}).Result()
	return result, err
}
