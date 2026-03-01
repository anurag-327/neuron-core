package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anurag-327/neuron-core/config"
	"github.com/anurag-327/neuron-core/conn"
	"github.com/anurag-327/neuron-core/engine"
	"github.com/anurag-327/neuron-core/pkg/logger"
	consoleLogger "github.com/anurag-327/neuron-core/pkg/logger/console"
	"github.com/anurag-327/neuron-core/runner"
	"github.com/anurag-327/neuron-core/runner/docker/pool"
	httpTransport "github.com/anurag-327/neuron-core/transport/http"
	"github.com/joho/godotenv"
)

func main() {

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// load config
	config.Load()

	// set console logger as global logger
	consoleLogger := consoleLogger.NewConsoleLogger("core")
	logger.SetGlobalLogger(consoleLogger)
	appLogger := logger.GetGlobalLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get docker client
	client, err := conn.GetDockerClient()
	if err != nil {
		log.Fatal("Error getting docker client")
	}

	// warm up docker pools
	if err := pool.InitDockerPool(ctx, client); err != nil {
		appLogger.Error(time.Now(), "Pool warm-up failed", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Pool warm-up failed: %v", err)
	}

	// create runner service
	runner := runner.NewRunner(client)
	// create execution service
	executionService := engine.NewExecutionService(runner)
	// create http handler
	handler := httpTransport.NewHandler(executionService)
	// create http router
	router := httpTransport.NewRouter(handler)
	// create http server
	srv := &http.Server{
		Addr:           ":" + config.PORT,
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// start server
	go startServer(srv)

	// handler graceful shutdown
	if err := gracefulShutdown(srv); err != nil {
		appLogger.Error(time.Now(), "Forced shutdown", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Forced shutdown: %v", err)
	}
}

func startServer(srv *http.Server) {
	appLogger := logger.GetGlobalLogger()
	appLogger.Info(time.Now(), "Server starting", map[string]interface{}{
		"port": config.PORT,
	})

	// start server on PORT
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Error(time.Now(), "Server start failed", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Server start failed: %v", err)
	}
}

func gracefulShutdown(srv *http.Server) error {

	// capture shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit
	log.Println("Shutdown signal received")

	// shutdown server with timeout after 30 seconds to allow existing connections to finish
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
