package validators

import (
	"errors"
	"github.com/plinkplenk/test-vortex/internal/api/schemas"
	"github.com/plinkplenk/test-vortex/internal/orders"
)

var (
	ErrExchangeNameNotProvided = errors.New("exchange not provided")

	ErrPairNotProvided  = errors.New("pair not provided")
	ErrDepthNotProvided = errors.New("depth not provided")

	ErrClientNotProvided       = errors.New("client not provided")
	ErrOrderHistoryNotProvided = errors.New("order history not provided")
)

func ValidateOrderBook(ob schemas.OrderBookCreate) error {
	if len(ob.ExchangeName) == 0 {
		return ErrExchangeNameNotProvided
	}
	if len(ob.Pair) == 0 {
		return ErrPairNotProvided
	}
	if len(ob.Depth) == 0 {
		return ErrDepthNotProvided
	}
	return nil
}

func ValidateClientHistory(ch schemas.ClientHistoryCreate) error {
	zeroClient := orders.Client{}
	zeroOrderHistory := orders.History{}
	if ch.Client == zeroClient {
		return ErrClientNotProvided
	}
	if ch.OrderHistory == zeroOrderHistory {
		return ErrOrderHistoryNotProvided
	}
	return nil
}
