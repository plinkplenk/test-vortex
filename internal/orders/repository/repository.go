package repository

import (
	"context"
	"github.com/plinkplenk/test-vortex/internal/orders"
)

type Repository interface {
	GetOrderBook(ctx context.Context, exchangeName, pair string) ([]orders.Depth, error)
	CreateOrderBook(ctx context.Context, exchangeName, pair string, orderBook []orders.Depth) error
	GetOrderHistory(ctx context.Context, client orders.Client) ([]*orders.History, error)
	CreateOrder(ctx context.Context, client orders.Client, order *orders.History) error
}
