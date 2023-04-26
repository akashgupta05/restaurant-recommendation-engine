package main

import (
	"sort"
	"time"

	"golang.org/x/exp/slices"
)

type sorter []func(agg *Aggregate) func(availableRestaurants *[]Restaurant)

var restaurantSorters = sorter{
	sortByFeaturedRestaurantsFunc,
	sortByPrimaryCuisinePrimaryCostFunc,
	sortByPrimaryCuisineSecondaryCostFunc,
	sortBySecondaryCuisinePrimaryCostFunc,
	sortByNewRestaurantsByRatingsFunc,
	sortByPrimaryCuisinePrimaryCostWithLessRatingFunc,
	sortByPrimaryCuisineSecondaryCostWithLessRatingFunc,
	sortBySecondaryCuisinePrimaryCostWithLessRatingFunc,
	sortByAnyCuisineAnyCostFunc,
}

func sortByFeaturedRestaurantsFunc(agg *Aggregate) func(availableRestaurants *[]Restaurant) {
	return agg.sortByFeaturedRestaurants
}

func sortByPrimaryCuisinePrimaryCostFunc(agg *Aggregate) func(availableRestaurants *[]Restaurant) {
	return agg.sortByPrimaryCuisinePrimaryCost
}

func sortByPrimaryCuisineSecondaryCostFunc(agg *Aggregate) func(availableRestaurants *[]Restaurant) {
	return agg.sortByPrimaryCuisineSecondaryCost
}

func sortBySecondaryCuisinePrimaryCostFunc(agg *Aggregate) func(availableRestaurants *[]Restaurant) {
	return agg.sortBySecondaryCuisinePrimaryCost
}

func sortByNewRestaurantsByRatingsFunc(agg *Aggregate) func(availableRestaurants *[]Restaurant) {
	return agg.sortByNewRestaurantsByRatings
}

func sortByPrimaryCuisinePrimaryCostWithLessRatingFunc(agg *Aggregate) func(availableRestaurants *[]Restaurant) {
	return agg.sortByPrimaryCuisinePrimaryCostWithLessRating
}

func sortByPrimaryCuisineSecondaryCostWithLessRatingFunc(agg *Aggregate) func(availableRestaurants *[]Restaurant) {
	return agg.sortByPrimaryCuisineSecondaryCostWithLessRating
}

func sortBySecondaryCuisinePrimaryCostWithLessRatingFunc(agg *Aggregate) func(availableRestaurants *[]Restaurant) {
	return agg.sortBySecondaryCuisinePrimaryCostWithLessRating
}

func sortByAnyCuisineAnyCostFunc(agg *Aggregate) func(availableRestaurants *[]Restaurant) {
	return agg.sortByAnyCuisineAnyCost
}

func (agg *Aggregate) sortByFeaturedRestaurants(availableRestaurants *[]Restaurant) {
	for _, restaurant := range *availableRestaurants {
		if restaurant.isRecommended && restaurant.cuisine == agg.primaryCuisine && restaurant.costBracket == agg.primaryCost {
			agg.setRestaurantId(restaurant.restaurantId)
			if ok := agg.circuitBreaker(); ok {
				return
			}
		}
	}
	if len(agg.restaurantIds) > 0 {
		return
	}

	for _, restaurant := range *availableRestaurants {
		if restaurant.isRecommended &&
			(isPrimaryCuisineSecondaryCosts(agg, restaurant.cuisine, restaurant.costBracket) ||
				isSecondaryCuisinesPrimaryCost(agg, restaurant.cuisine, restaurant.costBracket)) {
			agg.setRestaurantId(restaurant.restaurantId)
			if ok := agg.circuitBreaker(); ok {
				return
			}
		}
	}
}

func (agg *Aggregate) sortByNewRestaurantsByRatings(availableRestaurants *[]Restaurant) {
	sortFunc := func(i, j int) bool {
		if (*availableRestaurants)[i].rating == (*availableRestaurants)[j].rating {
			return (*availableRestaurants)[i].onboardedTime.After((*availableRestaurants)[j].onboardedTime)
		}
		return (*availableRestaurants)[i].rating > (*availableRestaurants)[j].rating
	}

	sort.Slice(*availableRestaurants, sortFunc)

	topFourRestaurantsCounter := 0
	now := time.Now()
	lastNHours := now.Add(time.Duration(NEW_RESTAURANT_TIME_IN_HOURS) * time.Hour)
	for _, restaurant := range *availableRestaurants {
		switch {
		case restaurant.onboardedTime.After(lastNHours):
			continue
		case topFourRestaurantsCounter >= 4:
			break
		default:
			if ok := agg.setRestaurantId(restaurant.restaurantId); ok {
				topFourRestaurantsCounter += 1
			}

			if ok := agg.circuitBreaker(); ok {
				return
			}
		}
	}
}

func (agg *Aggregate) sortByPrimaryCuisinePrimaryCost(availableRestaurants *[]Restaurant) {
	for _, restaurant := range *availableRestaurants {
		if restaurant.cuisine == agg.primaryCuisine && restaurant.costBracket == agg.primaryCost && restaurant.rating >= SECONDARY_HIGH_RATING {

			agg.setRestaurantId(restaurant.restaurantId)
			if ok := agg.circuitBreaker(); ok {
				return
			}
		}
	}
}

func (agg *Aggregate) sortByPrimaryCuisinePrimaryCostWithLessRating(availableRestaurants *[]Restaurant) {
	for _, restaurant := range *availableRestaurants {
		if restaurant.cuisine == agg.primaryCuisine && restaurant.costBracket == agg.primaryCost && restaurant.rating < SECONDARY_HIGH_RATING {

			agg.setRestaurantId(restaurant.restaurantId)
			if ok := agg.circuitBreaker(); ok {
				return
			}
		}
	}
}

func (agg *Aggregate) sortByPrimaryCuisineSecondaryCost(availableRestaurants *[]Restaurant) {
	for _, restaurant := range *availableRestaurants {
		if isPrimaryCuisineSecondaryCosts(agg, restaurant.cuisine, restaurant.costBracket) && restaurant.rating >= HIGH_RATING {

			agg.setRestaurantId(restaurant.restaurantId)
			if ok := agg.circuitBreaker(); ok {
				return
			}
		}
	}
}

func (agg *Aggregate) sortByPrimaryCuisineSecondaryCostWithLessRating(availableRestaurants *[]Restaurant) {
	for _, restaurant := range *availableRestaurants {
		if isPrimaryCuisineSecondaryCosts(agg, restaurant.cuisine, restaurant.costBracket) && restaurant.rating < HIGH_RATING {

			agg.setRestaurantId(restaurant.restaurantId)
			if ok := agg.circuitBreaker(); ok {
				return
			}
		}
	}
}

func (agg *Aggregate) sortBySecondaryCuisinePrimaryCost(availableRestaurants *[]Restaurant) {
	for _, restaurant := range *availableRestaurants {
		if isSecondaryCuisinesPrimaryCost(agg, restaurant.cuisine, restaurant.costBracket) && restaurant.rating >= HIGH_RATING {

			agg.setRestaurantId(restaurant.restaurantId)
			if ok := agg.circuitBreaker(); ok {
				return
			}
		}
	}
}

func (agg *Aggregate) sortBySecondaryCuisinePrimaryCostWithLessRating(availableRestaurants *[]Restaurant) {
	for _, restaurant := range *availableRestaurants {
		if isSecondaryCuisinesPrimaryCost(agg, restaurant.cuisine, restaurant.costBracket) && restaurant.rating < HIGH_RATING {

			agg.setRestaurantId(restaurant.restaurantId)
			if ok := agg.circuitBreaker(); ok {
				return
			}
		}
	}
}

func (agg *Aggregate) sortByAnyCuisineAnyCost(availableRestaurants *[]Restaurant) {
	for _, restaurant := range *availableRestaurants {

		agg.setRestaurantId(restaurant.restaurantId)
		if ok := agg.circuitBreaker(); ok {
			return
		}
	}
}

func isPrimaryCuisineSecondaryCosts(agg *Aggregate, cuisine Cuisine, costBracket int) bool {
	return cuisine == agg.primaryCuisine && slices.Contains(agg.secondaryCosts, costBracket)
}

func isSecondaryCuisinesPrimaryCost(agg *Aggregate, cuisine Cuisine, costBracket int) bool {
	return costBracket == agg.primaryCost && slices.Contains(agg.secondaryCuisines, cuisine)
}
