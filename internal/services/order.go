package services

import (
	"errors"
	"re-partners-challenge/internal/clients"
	"re-partners-challenge/internal/constants"
	"re-partners-challenge/internal/models"
	"sort"
)

type OrderI interface {
	CalculateOrderPackQty(order models.Order) ([]models.OrderPacks, error)
}

type Order struct {
	cache *clients.Cache
}

func NewOrderService(cache *clients.Cache) OrderI {
	return &Order{cache}
}

func (o *Order) CalculateOrderPackQty(order models.Order) ([]models.OrderPacks, error) {
	// If the number of items is 0 or negative, there is nothing to calculate
	if order.Items <= 0 {
		return nil, nil
	}

	// Gets the pack sizes available in the cache in memory
	packSizes, found := o.cache.Get(constants.PackSizesCacheKey)
	if !found {
		return nil, errors.New("pack sizes not found")
	}
	packSizesInt := packSizes.([]int)

	// Sorts the pack sizes from the largest to the smallest
	sort.Sort(sort.Reverse(sort.IntSlice(packSizesInt)))

	var packs []models.OrderPacks

	// For each available pack size
	for _, size := range packSizesInt {
		// If there are no more items to pack, stop the loop
		if order.Items <= 0 {
			break
		}

		// Calculates how many packs of this size are needed
		count := order.Items / size
		if count > 0 {
			// Add the pack with the calculated quantity
			packs = append(packs, models.OrderPacks{Size: size, Count: count})
			// Updates the amount of remaining items by subtracting the items we have already accounted
			order.Items %= size
		}
	}

	// If there are still items left, add packs of the smallest size required
	if order.Items > 0 {
		// Get the smallest pack size available
		smallestSize := packSizesInt[len(packSizesInt)-1]
		// Calculates how many packs of this size are needed for the remaining items
		count := (order.Items + smallestSize - 1) / smallestSize

		// Search if there is already a pack with the lowest value and just add the number of packs
		packFound := false
		for i, p := range packs {
			if p.Size == smallestSize {
				packs[i].Count += count
				packFound = true
				break
			}
		}

		// If there is no smaller pack, add a new pack already with the required quantity
		if !packFound {
			packs = append(packs, models.OrderPacks{Size: smallestSize, Count: count})
		}
	}

	// Return calculate packs
	return packs, nil
}
