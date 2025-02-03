package farm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kl09/seed-farm/internal/domain"
)

type Farmer struct {
	balancesRepository domain.BalancesRepository
	walletGenFn        func() (domain.Wallet, error)
	notifier           domain.Notifier
	goroutinesNum      int
}

func NewFarmer(
	balancesRepository domain.BalancesRepository,
	notifier domain.Notifier,
	walletGenFn func() (domain.Wallet, error),
	goroutinesNum int,
) *Farmer {
	return &Farmer{
		balancesRepository: balancesRepository,
		notifier:           notifier,
		walletGenFn:        walletGenFn,
		goroutinesNum:      goroutinesNum,
	}
}

func (e *Farmer) Run(ctx context.Context) {
	var (
		counter   int64
		wg        = sync.WaitGroup{}
		statsMx   = sync.Mutex{}
		startedAt = time.Now()
		now       = time.Now()
	)

	for i := 0; i < e.goroutinesNum; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()

			logger := slog.With("worker", worker)
			logger.Info("worker started")
			for {
				select {
				case <-ctx.Done():
					return
				default:
					wallet, err := e.walletGenFn()
					if err != nil {
						logger.Error(fmt.Sprintf("wallet generation: %s", err))
						continue
					}

					exists, err := e.balancesRepository.Exists(ctx, wallet.ETHAddressFormated())
					if err != nil {
						if errors.Is(err, context.Canceled) {
							return
						}
						logger.Error(fmt.Sprintf("address can't be checked: %s %s", err, wallet.String()))
					}
					if exists {
						err = e.notifier.WalletFound(ctx, wallet)
						if err != nil {
							logger.Error(fmt.Sprintf("wallet found notify error: %s %s", err, wallet.String()))
						}
					}

					atomic.AddInt64(&counter, 1)
					if time.Now().Add(-30 * time.Second).After(now) {
						if statsMx.TryLock() {
							c := atomic.LoadInt64(&counter)
							logger.Info(
								"report",
								"counter", c,
								"addresses/sec", c/int64(time.Since(startedAt).Seconds()),
							)
							now = time.Now()
							statsMx.Unlock()
						}
					}
				}
			}
		}(i)
	}

	wg.Wait()
}
