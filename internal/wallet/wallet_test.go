package wallet

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generateRandomMnemonic(t *testing.T) {
	s, err := generateRandomMnemonic(DefaultMnemonicBits)
	require.NoError(t, err)
	require.NotEmpty(t, s)
	require.Len(t, strings.Split(s, " "), 12)
}

func Test(t *testing.T) {
	w, err := NewWallet()
	require.NoError(t, err)
	require.Len(t, strings.Split(w.Mnemonic, " "), 12)
	require.Equal(t, "0x", w.ETHAddress[:2])
	require.NotEmpty(t, w.PrivateKey)
}
