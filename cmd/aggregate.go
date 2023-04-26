package main

import "sort"

type Aggregate struct {
	primaryCuisine    Cuisine
	secondaryCuisines []Cuisine
	primaryCost       int
	secondaryCosts    []int
	restaurantIdsMap  map[string]struct{}
	restaurantIds     []string
	maxRestaurantIds  int
}

func (a *Aggregate) setRestaurantId(id string) bool {
	if _, ok := a.restaurantIdsMap[id]; !ok {
		a.restaurantIdsMap[id] = struct{}{}
		a.restaurantIds = append(a.restaurantIds, id)
		return true
	}

	return false
}

func (a *Aggregate) circuitBreaker() bool {
	if len(a.restaurantIds) >= a.maxRestaurantIds {
		return true
	}

	return false
}

// setCuisines sets the primary cuisine and secondary cuisines
func (a *Aggregate) setCuisines(cuisineTrackings []CuisineTracking) {
	sort.Slice(cuisineTrackings, func(i, j int) bool {
		return cuisineTrackings[i].noOfOrders > cuisineTrackings[j].noOfOrders
	})

	secondaryCuisinesMap := map[Cuisine]struct{}{}
	length := len(cuisineTrackings)
	a.primaryCuisine = cuisineTrackings[length-1].Type

	for _, cuisineTracking := range cuisineTrackings {
		switch {
		case a.primaryCuisine == cuisineTracking.Type:
			continue
		case len(secondaryCuisinesMap) >= SECONDARY_CUISINES_LIMIT:
			break
		default:
			secondaryCuisinesMap[cuisineTracking.Type] = struct{}{}
		}
	}

	secondaryCuisines := []Cuisine{}
	for cuisine := range secondaryCuisinesMap {
		secondaryCuisines = append(secondaryCuisines, cuisine)
	}

	a.secondaryCuisines = secondaryCuisines
}

// setCosts sets the primary cost and secondary costs
func (a *Aggregate) setCosts(costTrackings []CostTracking) {
	sort.Slice(costTrackings, func(i, j int) bool {
		return costTrackings[i].noOfOrders > costTrackings[j].noOfOrders
	})

	secondaryCostMap := map[int]struct{}{}
	length := len(costTrackings)
	a.primaryCost = costTrackings[length-1].Type

	for _, costTracking := range costTrackings {
		switch {
		case a.primaryCost == costTracking.Type:
			continue
		case len(secondaryCostMap) >= SECONDARY_COSTS_LIMIT:
			break
		default:
			secondaryCostMap[costTracking.Type] = struct{}{}
		}
	}

	secondaryCosts := []int{}
	for cuisine := range secondaryCostMap {
		secondaryCosts = append(secondaryCosts, cuisine)
	}

	a.secondaryCosts = secondaryCosts
}
