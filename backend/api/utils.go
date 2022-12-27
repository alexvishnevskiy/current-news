package api

import (
	"fmt"
	"github.com/alexvishnevskiy/current-news/api/config"
	"github.com/alexvishnevskiy/current-news/api/extract"
	"strconv"
	"strings"
	"time"
)

type TimeStamp struct {
	Time  string
	Value int
}

func UpdateSortedSet(db *RedisDB, c config.Config) error {
	statsDict := extract.FetchAllData(c)
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

func UpdateHeadlines(db *RedisDB, c config.ConfigHead) error {
	articles, err := extract.FetchHead(c.Url, c.ApiKey, c.Sources, c.PageSize)
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

func GetHeadlines(db *RedisDB, key string) ([]extract.Article, error) {
	headlines, err := db.GetSet(key)
	if err != nil {
		return []extract.Article{}, err
	}
	return headlines, nil
}

func addToArchive(
	db *TimeseriesDB,
	c config.ConfigArchive,
	startTimestamp int64,
	startYear int,
	startMonth int,
	endYear int,
	endMonth int,
) error {

	// value to be added in start date
	addMonth, addYear := 0, 0
	switch c.Frequency {
	case "month":
		addMonth = 1
	case "year":
		addYear = 1
	}

	// TODO: parallelize it
	// iterate until we reach the end date
	for startYear <= endYear && startMonth <= endMonth {
		res, err := extract.FetchArchive(c.Url, c.ApiKey, startYear, startMonth)
		if err != nil {
			return err
		}

		err = db.Add(c.SetKey, startTimestamp, float64(res))
		if err != nil {
			return err
		}
		// update timestamp and date for next iteration
		startMonth, startYear = startMonth+addMonth, startYear+addYear
		startTimestamp++
		if startMonth > 12 {
			startMonth, startYear = 1, startYear+1
		}
	}
	return nil
}

func InitArchive(db *TimeseriesDB, c config.ConfigArchive) error {
	date := strings.Split(c.StartDate, "-")
	startYear, _ := strconv.Atoi(date[0])
	startMonth, _ := strconv.Atoi(date[1])
	endYear, endMonth, _ := time.Now().Date()

	// we only need data until this month
	// update function will update for the current month
	if int(endMonth)-1 == 0 {
		endYear, endMonth = endYear-1, 12
	} else {
		endMonth -= 1
	}

	err := addToArchive(db, c, 0, startYear, startMonth, endYear, int(endMonth))
	return err
}

func UpdateArchive(db *TimeseriesDB, c config.ConfigArchive) error {
	lastStamp, err := db.GetLast(c.SetKey)
	if err != nil {
		return err
	}
	// get dates to fetch data
	Year, Month, _ := time.Now().Date()
	err = addToArchive(db, c, lastStamp.Timestamp+1, Year, int(Month), Year, int(Month))
	return err
}

// return dict date -> value
func GetArchive(db *TimeseriesDB, c config.ConfigArchive) ([]TimeStamp, error) {
	lastStamp, err := db.GetLast(c.SetKey)
	if err != nil {
		return []TimeStamp{}, err
	}

	date := strings.Split(c.StartDate, "-")
	Year, _ := strconv.Atoi(date[0])
	Month, _ := strconv.Atoi(date[1])
	// value to be added in start date
	addMonth, addYear := 0, 0
	switch c.Frequency {
	case "month":
		addMonth = 1
	case "year":
		addYear = 1
	}

	// result dict
	res, err := db.GetSeries(c.SetKey, 0, lastStamp.Timestamp)
	if err != nil {
		return []TimeStamp{}, err
	}
	// more comfortable format
	timeseries := make([]TimeStamp, len(res))
	for i, stamp := range res {
		timeseries[i] = TimeStamp{fmt.Sprintf("%d-%d", Year, Month), int(stamp.Value)}
		Month, Year = Month+addMonth, Year+addYear
		if Month > 12 {
			Month, Year = 1, Year+1
		}
	}
	return timeseries, nil
}
