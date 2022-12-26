package api

import (
	"fmt"

	cat "github.com/alexvishnevskiy/current-news/api/categories"
	head "github.com/alexvishnevskiy/current-news/api/headlines"
)

func UpdateSortedSet(db *RedisDB, c cat.Config) error {
	statsDict := cat.FetchAllData(c)
	if len(statsDict) != 0 {
		for continent := range statsDict {
			for category := range statsDict[continent] {
				err := db.Add(continent, statsDict[continent][category], category)
				if err != nil {
					return err
				}
			}
		}
		fmt.Println("Data for categories is updated")
	} else {
		fmt.Println("Fetched data for categories is empty")
	}
	return nil
}

func GetCategoriesTop(db *RedisDB, continent string) ([]Member, error) {
	var topCategories []Member
	exists, err := db.SetExists(continent)
	if !exists || err != nil {
		return topCategories, err
	}

	topCategories, err = db.GetTop(continent)
	if err != nil {
		return topCategories, err
	}
	return topCategories, nil
}

func UpdateHeadlines(db *RedisDB, c head.ConfigAPI) error {
	articles, err := head.Fetch(c.Url, c.ApiKey, c.Sources, c.PageSize)
	if err != nil {
		return err
	}

	if len(articles) != 0 {
		for _, article := range articles {
			if err := db.AddSet(c.SetKey, article); err != nil {
				return err
			}
		}
		fmt.Println("Data for headlines is updated")
	} else {
		fmt.Println("Fetched data for headlines is empty")
	}
	return nil
}

func GetHeadlines(db *RedisDB, key string) ([]head.Article, error) {
	headlines, err := db.GetSet(key)
	if err != nil {
		return []head.Article{}, err
	}
	return headlines, nil
}
