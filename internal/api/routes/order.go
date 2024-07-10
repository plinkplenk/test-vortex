package routes

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/plinkplenk/test-vortex/internal/api/handlers"
	order "github.com/plinkplenk/test-vortex/internal/orders/service"
	"log/slog"
	"net/http"
)

func OrderRouter(ctx context.Context, orderService order.OrdersService, logger *slog.Logger) http.Handler {
	orderHandler := handlers.NewOrdersHandler(orderService, logger)
	r := chi.NewRouter()
	r.Get("/{exchange_name}/{pair}", orderHandler.GetOrderBook(ctx))
	r.Post("/", orderHandler.SaveOrderBook(ctx))
	r.Route(
		"/history", func(r chi.Router) {
			r.Get("/{client_name}/{exchange_name}", orderHandler.GetOrderHistory(ctx))
			r.Post("/", orderHandler.SaveOrder(ctx))
		},
	)
	return r
}
