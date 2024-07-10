package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/plinkplenk/test-vortex/internal/api/schemas"
	"github.com/plinkplenk/test-vortex/internal/api/validators"
	"github.com/plinkplenk/test-vortex/internal/orders"
	order "github.com/plinkplenk/test-vortex/internal/orders/service"
	"io"
	"log/slog"
	"net/http"
)

var (
	ErrOrderBookNotFound      = errors.New("order book not found")
	ErrLabelOrPairNotProvided = errors.New("label and pair not provided in query params")
)

type OrdersHandler struct {
	orderService order.OrdersService
	logger       *slog.Logger
}

func NewOrdersHandler(orderService order.OrdersService, logger *slog.Logger) *OrdersHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &OrdersHandler{orderService: orderService, logger: logger}
}

func (oh *OrdersHandler) GetOrderBook(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchangeName := chi.URLParam(r, "exchange_name")
		pair := chi.URLParam(r, "pair")
		depthOrders, err := oh.orderService.GetOrderBook(ctx, exchangeName, pair)
		if len(depthOrders) == 0 {
			if err != nil {
				oh.logger.Debug(
					"error while trying to get order book",
					"exchangeName", exchangeName,
					"pair", pair,
					"error", err,
				)
			}
			if err := response(
				j{"error": ErrOrderBookNotFound.Error()},
				http.StatusNotFound,
				w,
			); err != nil {
				logError(oh.logger, r, err)
			}
			return
		}
		if err := response(j{"depthOrders": depthOrders}, http.StatusOK, w); err != nil {
			logError(oh.logger, r, err)
			return
		}
	}
}

func (oh *OrdersHandler) SaveOrderBook(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logError(oh.logger, r, err)
			return
		}
		var orderToCreate schemas.OrderBookCreate
		if err := json.Unmarshal(body, &orderToCreate); err != nil {
			logError(oh.logger, r, err)
			if err := response(
				j{"error": "You must provide correct exchange name pair"},
				http.StatusBadRequest,
				w,
			); err != nil {
				logError(oh.logger, r, err)
			}
			return
		}
		if err := validators.ValidateOrderBook(orderToCreate); err != nil {
			if err := response(
				j{"error": err.Error()},
				http.StatusBadRequest,
				w,
			); err != nil {
				logError(oh.logger, r, err)
			}
			return
		}
		if err := oh.orderService.SaveOrderBook(
			ctx,
			orderToCreate.ExchangeName,
			orderToCreate.Pair,
			orderToCreate.Depth,
		); err != nil {
			oh.logger.Debug(
				"error while trying to save order book",
				"exchangeName", orderToCreate.ExchangeName,
				"pair", orderToCreate.Pair,
				"depth", orderToCreate.Depth,
				"error", err,
			)
			if err := response(j{"error": "something went wrong"}, http.StatusInternalServerError, w); err != nil {
				logError(oh.logger, r, err)
				return
			}
		}
		if err := response(nil, http.StatusCreated, w); err != nil {
			logError(oh.logger, r, err)
			return
		}
	}
}
func (oh *OrdersHandler) GetOrderHistory(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientName := chi.URLParam(r, "client_name")
		exchangeName := chi.URLParam(r, "exchange_name")
		label := r.URL.Query().Get("label")
		pair := r.URL.Query().Get("pair")
		if label == "" || pair == "" {
			if err := response(
				j{"error": ErrLabelOrPairNotProvided.Error()},
				http.StatusBadRequest,
				w,
			); err != nil {
				logError(oh.logger, r, err)
			}
			return
		}
		var client = orders.Client{
			ClientName:   clientName,
			ExchangeName: exchangeName,
			Label:        label,
			Pair:         pair,
		}
		orderHistory, err := oh.orderService.GetOrderHistory(ctx, client)
		if err != nil {
			logError(oh.logger, r, err)
			return
		}
		if orderHistory == nil {
			oh.logger.Debug(
				"error while trying to get order book",
				"exchangeName", exchangeName,
				"pair", pair,
				"error", err,
			)
			if err := response(
				j{"error": fmt.Sprintf("Order history with lable[%s] and pair[%s] not found", label, pair)},
				http.StatusNotFound,
				w,
			); err != nil {
				logError(oh.logger, r, err)
			}
			return
		}
		if err := response(j{"orderHistory": orderHistory}, http.StatusOK, w); err != nil {
			logError(oh.logger, r, err)
		}
	}
}
func (oh *OrdersHandler) SaveOrder(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logError(oh.logger, r, err)
			return
		}
		var clientHistory schemas.ClientHistoryCreate
		if err := json.Unmarshal(body, &clientHistory); err != nil {
			logError(oh.logger, r, err)
			if err := response(
				j{"error": "Client or History not provided"},
				http.StatusBadRequest,
				w,
			); err != nil {
				logError(oh.logger, r, err)
			}
			return
		}
		if err := validators.ValidateClientHistory(clientHistory); err != nil {
			if err := response(
				j{"error": err.Error()},
				http.StatusBadRequest,
				w,
			); err != nil {
				logError(oh.logger, r, err)
			}
			return
		}

		if err := oh.orderService.SaveOrder(ctx, clientHistory.Client, &clientHistory.OrderHistory); err != nil {
			oh.logger.Debug(
				"error while trying to save order book",
				"client", clientHistory.Client,
				"orderHistory", clientHistory.OrderHistory,
				"error", err,
			)
			if err := response(j{"error": "something went wrong"}, http.StatusInternalServerError, w); err != nil {
				logError(oh.logger, r, err)
			}
			return
		}
		if err := response(nil, http.StatusCreated, w); err != nil {
			logError(oh.logger, r, err)
			return
		}
	}
}
