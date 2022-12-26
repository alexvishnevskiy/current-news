package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexvishnevskiy/current-news/api/extract"
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

func (db *RedisDB) RemoveSet(key string) error {
	_, err := db.client.Del(db.Ctx, key).Result()
	return err
}

func (db *RedisDB) GetTop(continent string) ([]Member, error) {
	var inf int64 = 10000000
	leaderboard := make([]Member, 0)
	result, err := db.client.ZRangeWithScores(db.Ctx, continent, 0, inf).Result()

	if err != nil {
		return leaderboard, err
	}

	for i := len(result) - 1; i >= 0; i-- {
		leaderboard = append(leaderboard, Member{Member: fmt.Sprintf("%v", result[i].Member), Score: int(result[i].Score)})
	}
	return leaderboard, nil
}

func (db *RedisDB) AddSet(key string, value extract.Article) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = db.client.SAdd(db.Ctx, key, bytes).Result()
	return err
}

func (db *RedisDB) GetSet(key string) ([]extract.Article, error) {
	res, err := db.client.SMembers(db.Ctx, key).Result()
	if err != nil {
		return []extract.Article{}, err
	}

	// read and unmarshall members of the set
	setMembers := make([]extract.Article, len(res))
	for i, member := range res {
		var v extract.Article
		if err := json.Unmarshal([]byte(member), &v); err != nil {
			return []extract.Article{}, err
		}
		setMembers[i] = v
	}
	return setMembers, nil
}

func (db *RedisDB) Size(key string) (int64, error) {
	size, err := db.client.SCard(db.Ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return size, nil
}
