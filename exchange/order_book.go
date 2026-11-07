package exchange

import (
	"sort"
	"time"
)

type OrderBook struct {
	Bids []*PriceLevel // Sorted Descending (Highest buy prices first)
	Asks []*PriceLevel // Sorted Ascending (Lowest sell prices first)
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: make([]*PriceLevel, 0),
		Asks: make([]*PriceLevel, 0),
	}
}

// ProcessOrder accepts an aggressive inbound order and immediately scans the opposite book ladder side.
func (ob *OrderBook) ProcessOrder(takerOrder *Order) []*ExecutionReport {
	takerOrder.Timestamp = time.Now()
	var matches []*ExecutionReport

	if takerOrder.Side == Buy {
		matches = ob.matchOrder(takerOrder, &ob.Asks)
		// If taker order still has size remaining after eating liquidity, park it as a passive maker bid
		if takerOrder.Remaining() > 0 {
			ob.queueOrder(takerOrder, &ob.Bids, true)
		}
	} else {
		matches = ob.matchOrder(takerOrder, &ob.Bids)
		if takerOrder.Remaining() > 0 {
			ob.queueOrder(takerOrder, &ob.Asks, false)
		}
	}

	return matches
}

func (ob *OrderBook) matchOrder(taker *Order, makerLevels *[]*PriceLevel) []*ExecutionReport {
	var reports []*ExecutionReport

	for len(*makerLevels) > 0 && taker.Remaining() > 0 {
		level := (*makerLevels)[0]

		// Check if prices cross: Buy Taker must be >= Sell Maker | Sell Taker must be <= Buy Maker
		if (taker.Side == Buy && taker.Price < level.Price) || (taker.Side == Sell && taker.Price > level.Price) {
			break // No crossing match possible, halt matching sweep loop
		}

		// Loop through the FIFO time priority linked list on this specific price level
		for level.Head != nil && taker.Remaining() > 0 {
			maker := level.Head
			matchSize := taker.Remaining()
			if maker.Remaining() < matchSize {
				matchSize = maker.Remaining()
			}

			// Execute matching array adjustments
			taker.FilledQty += matchSize
			maker.FilledQty += matchSize
			level.TotalVol -= matchSize

			reports = append(reports, &ExecutionReport{
				MakerOrderID: maker.ID,
				TakerOrderID: taker.ID,
				Price:        level.Price, // Trade executes at the resting maker's price
				Amount:       matchSize,
				Timestamp:    time.Now(),
			})

			// If the resting maker is fully consumed, remove it from the linked list queue head
			if maker.Remaining() == 0 {
				level.Remove(maker)
			}
		}

		// If the entire price tier level is emptied out, prune the level index out of the ladder slice
		if level.TotalVol == 0 {
			*makerLevels = (*makerLevels)[1:]
		}
	}

	return reports
}

func (ob *OrderBook) queueOrder(order *Order, levels *[]*PriceLevel, descending bool) {
	// Look to see if an active price level struct already exists for this price point
	for _, level := range *levels {
		if level.Price == order.Price {
			level.Append(order)
			return
		}
	}

	// Instantiate a brand new price point container layer if it's the first order at this price
	newLevel := NewPriceLevel(order.Price)
	newLevel.Append(order)
	*levels = append(*levels, newLevel)

	// Keep our ladder array slices perfectly prioritized using native sort functions
	if descending {
		sort.Slice(*levels, func(i, j int) bool { return (*levels)[i].Price > (*levels)[j].Price })
	} else {
		sort.Slice(*levels, func(i, j int) bool { return (*levels)[i].Price < (*levels)[j].Price })
	}
}
