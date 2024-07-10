package repository

import (
	"context"
	"errors"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/plinkplenk/test-vortex/internal/orders"
)

var ErrOrderNotProvided = errors.New("order not provided")

type clickHouseRepository struct {
	db clickhouse.Conn
}

func NewClickHouseRepository(db clickhouse.Conn) Repository {
	return clickHouseRepository{
		db: db,
	}
}

func (r clickHouseRepository) GetOrderBook(ctx context.Context, exchangeName, pair string) (
	[]orders.Depth, error,
) {
	query := `SELECT asks FROM order_book WHERE exchange = ? AND pair = ?`
	rows, err := r.db.Query(ctx, query, exchangeName, pair)
	if err != nil {
		return nil, err
	}
	var asks []orders.Depth
	for rows.Next() {
		// asks is Tuple of two floats in clickhouse table
		var ask []float64
		if err := rows.Scan(&ask); err != nil {
			return nil, err
		}
		asks = append(
			asks, orders.Depth{
				Price:   ask[0],
				BaseQty: ask[1],
			},
		)
	}
	return asks, nil
}

func (r clickHouseRepository) CreateOrderBook(
	ctx context.Context, exchangeName, pair string, orderBook []orders.Depth,
) error {
	batch, err := r.db.PrepareBatch(ctx, "INSERT INTO order_book (exchange, pair, asks)")
	if err != nil {
		return err
	}
	for _, order := range orderBook {
		if err = batch.Append(
			exchangeName,
			pair,
			[]float64{order.Price, order.BaseQty},
		); err != nil {
			return err
		}
	}
	return batch.Send()
}

// GetOrderHistory returns history of order by client name, exchange name, label and pair
func (r clickHouseRepository) GetOrderHistory(ctx context.Context, client orders.Client) (
	[]*orders.History, error,
) {
	query := `
		SELECT side, type, base_qty, price, algorithm_name_placed, lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed
		FROM order_history 
		WHERE client_name = ? AND exchange_name = ? AND label = ? AND pair = ?`
	rows, err := r.db.Query(ctx, query, client.ClientName, client.ExchangeName, client.Label, client.Pair)
	if err != nil {
		return nil, err
	}
	var orderHistories []*orders.History
	for rows.Next() {
		orderHistory := orders.History{Client: client}
		if err := rows.Scan(
			&orderHistory.Side,
			&orderHistory.Type,
			&orderHistory.BaseQty,
			&orderHistory.Price,
			&orderHistory.AlgorithmNamePlaced,
			&orderHistory.LowestSellPrc,
			&orderHistory.HighestBuyPrc,
			&orderHistory.CommissionQuoteQty,
			&orderHistory.TimePlaced,
		); err != nil {
			return nil, err
		}
		orderHistories = append(orderHistories, &orderHistory)
	}
	return orderHistories, nil
}

func (r clickHouseRepository) CreateOrder(
	ctx context.Context, client orders.Client, order *orders.History,
) error {
	if order == nil {
		return ErrOrderNotProvided
	}
	query := `INSERT INTO order_history (
				client_name,
				exchange_name,
				label,
				pair,
				side,
				type,
				base_qty,
				price,
				algorithm_name_placed,
				lowest_sell_prc,
				highest_buy_prc,
				commission_quote_qty
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if err := r.db.Exec(
		ctx,
		query,
		client.ClientName,
		client.ExchangeName,
		client.Label,
		client.Pair,
		order.Side,
		order.Type,
		order.BaseQty,
		order.Price,
		order.AlgorithmNamePlaced,
		order.LowestSellPrc,
		order.HighestBuyPrc,
		order.CommissionQuoteQty,
	); err != nil {
		return err
	}
	return nil
}
