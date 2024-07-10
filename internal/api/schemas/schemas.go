package schemas

import "github.com/plinkplenk/test-vortex/internal/orders"

type ClientHistoryCreate struct {
	Client       orders.Client  `json:"client"`
	OrderHistory orders.History `json:"orderHistory"`
}

type OrderBookCreate struct {
	ExchangeName string         `json:"exchangeName"`
	Pair         string         `json:"pair"`
	Depth        []orders.Depth `json:"depth"`
}
