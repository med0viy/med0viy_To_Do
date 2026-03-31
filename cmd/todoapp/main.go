package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_postgres_pool "github.com/med0viy/practika/internal/core/repository/postgres/pool"
	core_http_maddleware "github.com/med0viy/practika/internal/core/transport/http/middleware"
	core_http_server "github.com/med0viy/practika/internal/core/transport/http/server"
	users_postgres_repository "github.com/med0viy/practika/internal/features/users/repository/postgres"
	users_service "github.com/med0viy/practika/internal/features/users/service"
	users_transport_http "github.com/med0viy/practika/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initiazling postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()



	logger.Debug("initiazling feauture", zap.String("feature", "users"))
	userRepository := users_postgres_repository.NewUsersRepository(pool)
	userServise := users_service.NewUsersServise(userRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(userServise)

	logger.Debug("initiazling HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_maddleware.RequestID(),
		core_http_maddleware.Logger(logger),
		core_http_maddleware.Panic(),
		core_http_maddleware.Trace(),
	)
	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
