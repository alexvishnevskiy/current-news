package api

import "fmt"

func UpdateTable(db *RedisDB, c Config) error {
	statsDict := FetchAllData(c)
	if len(statsDict) != 0 {
		for continent := range statsDict {
			for category := range statsDict[continent] {
				err := db.Add(continent, statsDict[continent][category], category)
				if err != nil {
					return err
				}
			}
		}
		fmt.Println("Data is updated")
	} else {
		fmt.Println("Fetched data is empty")
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
