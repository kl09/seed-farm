package notifier

import (
	"context"
	"os"

	"github.com/kl09/seed-farm/internal/domain"
)

type FileNotifier struct{}

func (n *FileNotifier) WalletFound(ctx context.Context, w domain.Wallet) error {
	return os.WriteFile(w.ETHAddress+".txt", []byte(w.String()), 0644)
}
