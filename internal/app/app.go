package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/plinkplenk/test-vortex/internal/api/middleware"
	"github.com/plinkplenk/test-vortex/internal/api/routes"
	"github.com/plinkplenk/test-vortex/internal/config"
	ordersRepository "github.com/plinkplenk/test-vortex/internal/orders/repository"
	order "github.com/plinkplenk/test-vortex/internal/orders/service"
	"log"
	"log/slog"
	"net/http"
	"time"
)

func setupRouters(
	orderService order.OrdersService, logger *slog.Logger, middlewares ...middleware.Middleware,
) http.Handler {
	return routes.NewRouter(orderService, logger, middlewares...)
}

func connectToClickhouse(clickhouseCfg config.Clickhouse, debug bool) (clickhouse.Conn, error) {
	addr := fmt.Sprintf("%s:%s", clickhouseCfg.Host, clickhouseCfg.Port)
	conn, err := clickhouse.Open(
		&clickhouse.Options{
			Addr: []string{addr},
			Auth: clickhouse.Auth{
				Database: "default",
				Username: clickhouseCfg.User,
				Password: clickhouseCfg.Password,
			},
			Debug: debug,
			Debugf: func(format string, v ...any) {
				log.Printf(format+"\n", v...)
			},
			Settings: clickhouse.Settings{
				"max_execution_time": 60,
			},
			Compression: &clickhouse.Compression{
				Method: clickhouse.CompressionLZ4,
			},
			DialTimeout:          time.Second * 10,
			MaxOpenConns:         300,
			MaxIdleConns:         10,
			ConnMaxLifetime:      30 * time.Second,
			ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
			BlockBufferSize:      10,
			MaxCompressionBuffer: 10240,
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "order-service", Version: "0.1"},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type App struct {
	env    config.ENV
	debug  bool
	config config.Config
	dbConn clickhouse.Conn
	server *http.Server
	logger *slog.Logger
}

type Params struct {
	Config config.Config
	Logger *slog.Logger
	Debug  bool
}

func New(params Params) (*App, error) {
	chConn, err := connectToClickhouse(params.Config.Clickhouse, params.Debug)
	if err != nil {
		return nil, err
	}
	orderService := order.New(ordersRepository.NewClickHouseRepository(chConn), params.Config.Server.Timeout)
	loggerMiddleware := middleware.NewLoggerMiddleware(params.Logger)
	handler := setupRouters(orderService, params.Logger, loggerMiddleware.Log)
	server := setupServer(params.Config.Server.Port, handler)
	return &App{
		env:    params.Config.ENV,
		dbConn: chConn,
		logger: params.Logger,
		config: params.Config,
		server: server,
		debug:  params.Debug,
	}, nil
}

func setupServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: handler,
	}
}

func (a *App) Run() error {
	a.logger.Info("Running server", "address", a.server.Addr)
	return a.server.ListenAndServe()
}

// Stop
// gracefully shuts down app with 30 seconds time out
func (a *App) Stop() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	go func() {
		<-shutdownCtx.Done()
		if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
			log.Fatal("graceful shutdown timeout\nForcing exit")
		}
	}()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error while trying to shutdown server", "error", err)
	}
	if err := a.dbConn.Close(); err != nil {
		slog.Error("Error on db connection close", "error", err)
	}
}
