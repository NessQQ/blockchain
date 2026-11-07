package exchange

// PriceLevel groups a FIFO execution queue of orders sharing the exact same price tier.
type PriceLevel struct {
	Price     float64
	TotalVol  float64
	Head      *Order
	Tail      *Order
}

func NewPriceLevel(price float64) *PriceLevel {
	return &PriceLevel{Price: price}
}

// Append inserts a incoming order at the very end of the time priority list line: O(1)
func (pl *PriceLevel) Append(order *Order) {
	pl.TotalVol += order.Remaining()
	if pl.Head == nil {
		pl.Head = order
		pl.Tail = order
		return
	}
	pl.Tail.Next = order
	order.Prev = pl.Tail
	pl.Tail = order
}

// Remove extracts an order from anywhere in the chain if cancelled or fully filled: O(1)
func (pl *PriceLevel) Remove(order *Order) {
	pl.TotalVol -= order.Remaining()
	
	if order.Prev != nil {
		order.Prev.Next = order.Next
	} else {
		pl.Head = order.Next // If it was the head, advance head forward
	}
	
	if order.Next != nil {
		order.Next.Prev = order.Prev
	} else {
		pl.Tail = order.Prev // If it was the tail, pull tail backward
	}
	
	order.Next = nil
	order.Prev = nil
}
