package routes

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/plinkplenk/test-vortex/internal/api/middleware"
	order "github.com/plinkplenk/test-vortex/internal/orders/service"
	"log/slog"
	"net/http"
)

func NewRouter(
	orderService order.OrdersService, logger *slog.Logger, middlewares ...middleware.Middleware,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares...)
	r.Mount("/orders", OrderRouter(context.Background(), orderService, logger))
	return r
}
