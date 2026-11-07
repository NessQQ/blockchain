// Copyright (C) 2026 Lenar Khabibulin <khabibulin.lenar@gmail.com>
// Licensed under the GNU Affero General Public License v3 (AGPL-3.0)
package exchange

import "time"

type Side int

const (
	Buy Side = iota
	Sell
)

// Order represents an individual active limit order submitted by a trader.
type Order struct {
	ID        string    `json:"id"`
	WalletID  string    `json:"walletId"`
	Side      Side      `json:"side"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`    // Total initial size requested
	FilledQty float64   `json:"filledQty"` // Amount filled so far
	Timestamp time.Time `json:"timestamp"` // Crucial for Time Priority sorting rules
	Next      *Order    `json:"-"`         // Pointer for our $O(1)$ Price Doubly Linked List node
	Prev      *Order    `json:"-"`         // Pointer for our $O(1)$ Price Doubly Linked List node
}

// Remaining returns the unfilled residual quantity left on the order book.
func (o *Order) Remaining() float64 {
	return o.Amount - o.FilledQty
}

// ExecutionReport captures an immediate match event logged by the engine core.
type ExecutionReport struct {
	MakerOrderID string    `json:"makerOrderId"` // Passive order sitting in book
	TakerOrderID string    `json:"takerOrderId"` // Aggressive inbound hitting order
	Price        float64   `json:"price"`
	Amount       float64   `json:"amount"`
	Timestamp    time.Time `json:"timestamp"`
}
