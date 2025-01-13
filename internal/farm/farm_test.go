package farm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kl09/seed-farm/internal/wallet"
	"github.com/kl09/seed-farm/test"
	"go.uber.org/mock/gomock"
)

func TestFarmer_Run(t *testing.T) {
	t.Run("exists - notifier error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := test.NewMockBalancesRepository(ctrl)
		notifier := test.NewMockNotifier(ctrl)
		repo.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
		notifier.EXPECT().WalletFound(gomock.Any(), gomock.Any()).Return(errors.New("some error"))

		ctx, cancelFn := context.WithCancelCause(context.Background())
		farmer := NewFarmer(repo, notifier, wallet.NewWallet, 1)
		farmer.Run(ctx, cancelFn)
	})

	t.Run("exists - notifier success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := test.NewMockBalancesRepository(ctrl)
		notifier := test.NewMockNotifier(ctrl)
		repo.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
		notifier.EXPECT().WalletFound(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		ctx, _ := context.WithTimeout(context.Background(), time.Second/2)
		ctx, cancelFn := context.WithCancelCause(ctx)
		farmer := NewFarmer(repo, notifier, wallet.NewWallet, 1)
		farmer.Run(ctx, cancelFn)
	})

	t.Run("not exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := test.NewMockBalancesRepository(ctrl)
		notifier := test.NewMockNotifier(ctrl)
		repo.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()

		ctx, _ := context.WithTimeout(context.Background(), time.Second/2)
		ctx, cancelFn := context.WithCancelCause(ctx)
		farmer := NewFarmer(repo, notifier, wallet.NewWallet, 1)
		farmer.Run(ctx, cancelFn)
	})
}
