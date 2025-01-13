package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kl09/seed-farm/internal/farm"
	"github.com/kl09/seed-farm/internal/notifier"
	"github.com/kl09/seed-farm/internal/repository"
	"github.com/kl09/seed-farm/internal/wallet"
)

const (
	postgresDSN = "postgres://127.0.0.1:5432/wallet"
)

func main() {
	ctx, cancelFn := context.WithCancelCause(context.Background())

	db, err := dialPG(ctx, postgresDSN)
	if err != nil {
		slog.Error(fmt.Sprintf("dial PG: %s", err))
		os.Exit(1)
	}

	fNotifier := notifier.FileNotifier{}
	balancesRepo := repository.NewBalancesRepository(db)
	farmer := farm.NewFarmer(balancesRepo, &fNotifier, wallet.NewWallet, 5)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		slog.Info("signal received, shutting down")
		signal.Reset()
		cancelFn(errors.New("shutdown"))
	}()

	farmer.Run(ctx, cancelFn)
}

func dialPG(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	return pool, nil
}
