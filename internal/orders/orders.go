package orders

import (
	"github.com/google/uuid"
	"time"
)

type Depth struct {
	Price   float64 `json:"price"`
	BaseQty float64 `json:"baseQty"`
}
type Book struct {
	ID       uuid.UUID `json:"id"`
	Exchange string    `json:"exchange"`
	Pair     string    `json:"pair"`
	Asks     []Depth   `json:"asks"`
	Bids     []Depth   `json:"bids"`
}

type Client struct {
	ClientName   string `json:"clientName"`
	ExchangeName string `json:"exchangeName"`
	Label        string `json:"label"`
	Pair         string `json:"pair"`
}

type History struct {
	Client
	Side                string    `json:"side"`
	Type                string    `json:"type"`
	BaseQty             float64   `json:"baseQty"`
	Price               float64   `json:"price"`
	AlgorithmNamePlaced string    `json:"algorithmNamePlaced"`
	LowestSellPrc       float64   `json:"lowestSellPrc"`
	HighestBuyPrc       float64   `json:"highestBuyPrc"`
	CommissionQuoteQty  float64   `json:"commissionQuoteQty"`
	TimePlaced          time.Time `json:"timePlaced"`
}
