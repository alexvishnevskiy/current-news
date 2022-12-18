package api

import "fmt"

func UpdateTable(db *RedisDB) error {
	statsDict := FetchAllData()
	if len(statsDict) != 0 {
		for continent := range statsDict {
			for category := range statsDict[continent] {
				err := db.Add(continent, statsDict[continent][category], category)
				if err != nil {
					return err
				}
			}
		}
	} else {
		fmt.Println("Fetched data is empty")
	}
	return nil
}

func GetData(db *RedisDB, continent string) ([]string, error) {
	exists, err := db.SetExists(continent)
	if !exists || err != nil {
		return []string{}, err
	}

	topk, err := db.GetTop(continent)
	if err != nil {
		return []string{}, err
	}
	return topk, err
}
