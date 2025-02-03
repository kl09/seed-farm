package notifier

import (
	"context"

	"github.com/kl09/seed-farm/internal/domain"
)

type TGNotifier struct{}

func (n *TGNotifier) WalletFound(ctx context.Context, w domain.Wallet) error {
	panic("not implemented")
}
