//go:generate mockgen -typed -source=./domain.go -destination=../../test/mock_services_domain.go -package=test
package domain

import (
	"context"
	"fmt"
	"strings"
)

// Wallet is a struct that contains the information of a wallet.
type Wallet struct {
	ETHAddress string
	PrivateKey string
	Mnemonic   string
}

func (w Wallet) ETHAddressFormated() string {
	return strings.ToLower(w.ETHAddress)[2:]
}

func (w Wallet) String() string {
	return fmt.Sprintf(
		"ETHAddress: %s\nPrivateKey: %s\nMnemonic: %s\n",
		w.ETHAddress, w.PrivateKey, w.Mnemonic,
	)
}

type BalancesRepository interface {
	Exists(ctx context.Context, address string) (exists bool, err error)
}

type Notifier interface {
	WalletFound(ctx context.Context, wallet Wallet) error
}
