package api

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

type RedisDB struct {
	client *redis.Client
	Ctx    context.Context
}

type Member struct {
	Member string
	Score  int
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

func (db *RedisDB) SetExists(continent string) (bool, error) {
	res, err := db.client.Exists(db.Ctx, continent).Result()

	if err != nil {
		return false, err
	} else if res == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (db *RedisDB) MemberExists(continent string, category string) (bool, error) {
	_, err := db.client.ZScore(db.Ctx, continent, category).Result()

	switch err {
	case redis.Nil:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (db *RedisDB) Add(continent string, score int, category string) error {
	_, err := db.client.ZAdd(
		db.Ctx, continent, redis.Z{Score: float64(score), Member: category}).Result()
	return err
}

func (db *RedisDB) RemoveMember(continent string, category string) error {
	_, err := db.client.ZRem(db.Ctx, continent, category).Result()
	return err
}

func (db *RedisDB) RemoveSet(continent string) error {
	_, err := db.client.Del(db.Ctx, continent).Result()
	return err
}

func (db *RedisDB) GetTop(continent string) ([]Member, error) {
	var inf int64 = 10000000
	leaderboard := make([]Member, 0)
	result, err := db.client.ZRangeWithScores(db.Ctx, continent, 0, inf).Result()

	if err != nil {
		return leaderboard, err
	}

	for _, record := range result {
		leaderboard = append(leaderboard, Member{Member: fmt.Sprintf("%v", record.Member), Score: int(record.Score)})
	}
	return leaderboard, nil
}
