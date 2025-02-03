package farm

import (
	"context"
	"testing"
	"time"

	"github.com/kl09/seed-farm/internal/wallet"
	"github.com/kl09/seed-farm/test"
	"go.uber.org/mock/gomock"
)

func TestFarmer_Run(t *testing.T) {
	t.Run("exists - notifier success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := test.NewMockBalancesRepository(ctrl)
		notifier := test.NewMockNotifier(ctrl)
		var exists bool
		repo.EXPECT().Exists(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, address string) (bool, error) {
			if !exists {
				exists = true
				return true, nil
			}
			return false, nil
		}).AnyTimes()
		notifier.EXPECT().WalletFound(gomock.Any(), gomock.Any()).Return(nil)

		ctx, _ := context.WithTimeout(context.Background(), time.Second/2)
		farmer := NewFarmer(repo, notifier, wallet.NewWallet, 1)
		farmer.Run(ctx)
	})

	t.Run("not exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := test.NewMockBalancesRepository(ctrl)
		notifier := test.NewMockNotifier(ctrl)
		repo.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()

		ctx, _ := context.WithTimeout(context.Background(), time.Second/2)
		farmer := NewFarmer(repo, notifier, wallet.NewWallet, 1)
		farmer.Run(ctx)
	})
}
