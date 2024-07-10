package service

import (
	"context"
	"github.com/plinkplenk/test-vortex/internal/orders"
	ordersRepository "github.com/plinkplenk/test-vortex/internal/orders/repository"
	"time"
)

//go:generate mockgen -source=orders.go -destination=mocks/mock.go
type OrdersService interface {
	GetOrderBook(ctx context.Context, exchangeName, pair string) ([]orders.Depth, error)
	SaveOrderBook(ctx context.Context, exchangeName, pair string, orderBook []orders.Depth) error
	GetOrderHistory(ctx context.Context, client orders.Client) ([]*orders.History, error)
	SaveOrder(ctx context.Context, client orders.Client, order *orders.History) error
}

type orderService struct {
	repository ordersRepository.Repository
	timeout    time.Duration
}

func New(repository ordersRepository.Repository, timeout time.Duration) OrdersService {
	return orderService{
		repository: repository,
		timeout:    timeout,
	}
}

func (s orderService) GetOrderBook(ctx context.Context, exchangeName, pair string) ([]orders.Depth, error) {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	depth, err := s.repository.GetOrderBook(c, exchangeName, pair)
	if err != nil {
		return nil, err
	}
	return depth, nil
}

func (s orderService) SaveOrderBook(ctx context.Context, exchangeName, pair string, orderBook []orders.Depth) error {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	err := s.repository.CreateOrderBook(c, exchangeName, pair, orderBook)
	return err
}

func (s orderService) GetOrderHistory(ctx context.Context, client orders.Client) ([]*orders.History, error) {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	history, err := s.repository.GetOrderHistory(c, client)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (s orderService) SaveOrder(ctx context.Context, client orders.Client, order *orders.History) error {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	err := s.repository.CreateOrder(c, client, order)
	if err != nil {
		return err
	}
	return nil
}
