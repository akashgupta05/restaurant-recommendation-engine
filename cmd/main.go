package main

import (
	"fmt"
	"time"
)

const (
	SouthIndian Cuisine = iota
	NorthIndian
	Chinese
)

const SECONDARY_CUISINES_LIMIT int = 2
const SECONDARY_COSTS_LIMIT int = 2
const HIGH_RATING float64 = 4.5
const SECONDARY_HIGH_RATING float64 = 4.0
const NEW_RESTAURANT_TIME_IN_HOURS int64 = -48

type Cuisine int

type Restaurant struct {
	restaurantId  string
	cuisine       Cuisine
	costBracket   int
	rating        float64
	isRecommended bool
	onboardedTime time.Time
}

type CuisineTracking struct {
	Type       Cuisine
	noOfOrders int
}

type CostTracking struct {
	Type       int
	noOfOrders int
}

type User struct {
	cuisines    []CuisineTracking
	costBracket []CostTracking
}

func getRestaurantRecommendations(user User, availableRestaurants []Restaurant) []string {
	agg := Aggregate{
		restaurantIdsMap: map[string]struct{}{},
		restaurantIds:    []string{},
		maxRestaurantIds: 100,
	}
	agg.setCosts(user.costBracket)
	agg.setCuisines(user.cuisines)

	for _, sorter := range restaurantSorters {
		sorter(&agg)(user, &availableRestaurants)
		if ok := agg.circuitBreaker(); ok {
			break
		}
	}
	return agg.restaurantIds
}

func main() {
	user := User{
		cuisines: []CuisineTracking{
			{Type: SouthIndian, noOfOrders: 5},
			{Type: NorthIndian, noOfOrders: 10},
			{Type: Chinese, noOfOrders: 3},
		},
		costBracket: []CostTracking{
			{Type: 1, noOfOrders: 7},
			{Type: 2, noOfOrders: 12},
			{Type: 5, noOfOrders: 2},
		},
	}

	restaurants := []Restaurant{
		{
			restaurantId:  "1",
			cuisine:       NorthIndian,
			costBracket:   2,
			rating:        4.4,
			isRecommended: true,
			onboardedTime: time.Now().Add(-1 * time.Hour),
		},
		{
			restaurantId:  "3",
			cuisine:       NorthIndian,
			costBracket:   5,
			rating:        4.1,
			isRecommended: false,
			onboardedTime: time.Now().Add(-4 * time.Hour),
		},
		{
			restaurantId:  "2",
			cuisine:       NorthIndian,
			costBracket:   5,
			rating:        4.1,
			isRecommended: false,
			onboardedTime: time.Now().Add(-4 * time.Hour),
		},
	}

	restaurantIds := getRestaurantRecommendations(user, restaurants)
	fmt.Println("Restaurant Ids \n", restaurantIds)
}
