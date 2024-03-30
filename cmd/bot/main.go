package main

import (
	"context"
	"errors"
	"flag"
	"github.com/makarellav/iad-dialogflow-bot/internal/models"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type bot struct {
	logger *slog.Logger
	coin   *models.CoinModel
}

func main() {
	addr := flag.String("addr", ":8080", "Bot port")
	baseUrl := flag.String("base", "https://api.coincap.io/v2/assets", "Coin API URL")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	b := bot{
		logger: logger,
		coin:   &models.CoinModel{BaseURL: *baseUrl},
	}

	srv := http.Server{
		Addr:         *addr,
		Handler:      b.handlers(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	errCh := make(chan error)

	go func() {
		quitCh := make(chan os.Signal, 1)

		signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
		s := <-quitCh

		b.logger.Info("shutting down the server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		errCh <- srv.Shutdown(ctx)
	}()

	b.logger.Info("starting the server", "addr", *addr)
	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err.Error())
	}

	err = <-errCh

	if err != nil {
		log.Fatal(err.Error())
	}

	b.logger.Info("stopped the server", "addr", srv.Addr)
}
