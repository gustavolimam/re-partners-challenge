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
		return nil, errors.New("no items to pack")
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
	var pack, bestPack = make(map[int]int), make(map[int]int)
	remainingItems := order.Items
	shouldContinue := true

	for shouldContinue {
		for i, size := range packSizesInt {
			// The sizes are in descending order, so if the value of i is equal to 0,
			// it means that there is no more package size available without the number of items exceeding that requested.
			if i == 0 && len(pack) > 0 {
				// If no packs of this size exist, remove this from the options to avoid unnecessary loop
				if pack[size] == 0 && len(packSizesInt) > 1 {
					packSizesInt = packSizesInt[1:]
					break
				}

				// If you don't have a smaller package, we return it the way it is, because this is the best scenario
				if i == len(packSizesInt)-1 {
					shouldContinue = false
					pack[size]++
					remainingItems -= size
					bestPack = pack
					continue
				}

				// Removes the largest package you have and increases the amount of items left with the value of this package
				pack[size]--
				remainingItems += size
				continue
			}

			// Calculates how many packs of this size are needed
			packsToFill := remainingItems / size
			// Update the remainingItems
			remainingItems = remainingItems - (packsToFill * size)
			// Store the number of packs for this size.
			pack[size] += packsToFill

			if remainingItems == 0 {
				// create a copy of this pack map to avoid issues with pointer
				copyPack := copyPackMap(pack)

				// validate if all items was packed, if yes we choose this as the bestPack option and stop the loop
				if allItemsIsPacked(order.Items, copyPack) {
					shouldContinue = false
					bestPack = copyPack
				}
				break
			}

			// If this is the smallest pack size we add to fill all the items we need
			if remainingItems > 0 && i == len(packSizesInt)-1 {
				// create a copy of this pack map to avoid issues with pointer
				copyPack := copyPackMap(pack)
				copyPack[size]++

				// validate if all items was packed
				if allItemsIsPacked(order.Items, copyPack) {
					bestPack = copyPack
					shouldContinue = false
				}
				break
			}
		}
	}

	// Cast the map to OrderPack model
	for _, key := range packSizesInt {
		// if there's no pack with this size, just go to the next interaction
		if bestPack[key] == 0 {
			continue
		}
		packs = append(packs, models.OrderPacks{
			Size:  key,
			Count: bestPack[key],
		})
	}

	return packs, nil
}

// check if all items on this pack is equal to items on order
func allItemsIsPacked(totalItems int, orderPacks map[int]int) bool {
	total := int(0)
	for packSize, count := range orderPacks {
		total += packSize * count
	}
	return totalItems == total
}

// creates a copy of a pack map
func copyPackMap(pack map[int]int) map[int]int {
	packCopied := make(map[int]int)
	for key, value := range pack {
		packCopied[key] = value
	}
	return packCopied
}
