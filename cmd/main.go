package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"senioravanti.ru/internal/s3-service/bootstrap"
	"senioravanti.ru/internal/s3-service/rest/routers"
)

func main() {
	app, err := bootstrap.New(); if err != nil {
		slog.Error("failed to bootstrap application", err)
		os.Exit(1)
	}

	routers.SetUp(app)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			slog.Error("failure", 
			slog.String("error", err.Error()))
		}
	}()

	slog.Info("server successfully started")

	<-done
	slog.Info("shutdown the server ...")

	timeout := "3s"
	timeoutDuration, _ := time.ParseDuration(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown the server",
		slog.String("error", err.Error()))
		return
	}

	<-ctx.Done()
	slog.Info("server successfully stopped", slog.String("timeout", timeout))
}