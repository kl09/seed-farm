package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"unsafe"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kl09/seed-farm/internal/domain"
	"github.com/pkg/errors"
	"github.com/planxnx/ethereum-wallet-generator/bip39"
)

const (
	// DefaultMnemonicBits is the default number of bits to use when generating a mnemonic. default is 128 bits (12 words).
	DefaultMnemonicBits = 128
)

var (
	// DefaultBaseDerivationPath is the base path from which custom derivation endpoints
	// are incremented. As such, the first account will be at m/44'/60'/0'/0, the second
	// at m/44'/60'/0'/1, etc
	DefaultBaseDerivationPath = accounts.DefaultBaseDerivationPath
)

func NewWallet() (domain.Wallet, error) {
	mnemonic, err := generateRandomMnemonic(DefaultMnemonicBits)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("generate random mnemonic: %w", err)
	}

	privateKey, err := deriveWallet(
		bip39.NewSeed(mnemonic, ""),
		DefaultBaseDerivationPath,
	)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("derive wallet: %w", err)
	}

	privateKeyString, err := privateKeyToString(privateKey)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("private key to string: %w", err)
	}
	publicKeyString, err := publicKeyToString(&privateKey.PublicKey)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("public key to string: %w", err)
	}

	return domain.Wallet{
		ETHAddress: publicKeyString,
		PrivateKey: privateKeyString,
		Mnemonic:   mnemonic,
	}, nil
}

func NewWalletByMnemonic(mnemonic string) (domain.Wallet, error) {
	privateKey, err := deriveWallet(
		bip39.NewSeed(mnemonic, ""),
		DefaultBaseDerivationPath,
	)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("derive wallet: %w", err)
	}

	privateKeyString, err := privateKeyToString(privateKey)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("private key to string: %w", err)
	}
	publicKeyString, err := publicKeyToString(&privateKey.PublicKey)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("public key to string: %w", err)
	}

	return domain.Wallet{
		ETHAddress: publicKeyString,
		PrivateKey: privateKeyString,
		Mnemonic:   mnemonic,
	}, nil
}

func privateKeyToString(key *ecdsa.PrivateKey) (string, error) {
	if key == nil {
		return "", errors.New("private key is nil")
	}

	privateKeyBytes := crypto.FromECDSA(key)
	privateHex := make([]byte, len(privateKeyBytes)*2)
	hex.Encode(privateHex, privateKeyBytes)
	return b2s(privateHex), nil
}

func publicKeyToString(key *ecdsa.PublicKey) (string, error) {
	publicKeyBytes := crypto.Keccak256(crypto.FromECDSAPub(key)[1:])[12:]
	if len(publicKeyBytes) > common.AddressLength {
		publicKeyBytes = publicKeyBytes[len(publicKeyBytes)-common.AddressLength:]
	}
	pubHex := make([]byte, len(publicKeyBytes)*2+2)
	copy(pubHex[:2], "0x")
	hex.Encode(pubHex[2:], publicKeyBytes)
	pubString := b2s(pubHex)

	return pubString, nil
}

func deriveWallet(seed []byte, path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("new master node: %w", err)
	}

	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, fmt.Errorf("derive key: %w", err)
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("ec private key: %w", err)
	}

	return privateKey.ToECDSA(), nil
}

// generateRandomMnemonic returns a new random mnemonic(BIP39) with the given bit size.
func generateRandomMnemonic(bitSize int) (string, error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", fmt.Errorf("new entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("new mnemonic: %w", err)
	}

	return mnemonic, nil
}

// b2s converts a byte slice to a string without memory allocation.
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
