package mnemonic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FromMnemonicToShort_FromShortToMnemonic(t *testing.T) {
	mnemonic := "squeeze hill cube network mobile catalog plate yellow obtain oppose floor interest"
	short := FromMnemonicToShort(mnemonic)
	require.Equal(t, "1692.861.427.1190.1139.287.1329.2040.1220.1244.715.941.", short)
	got, err := FromShortToMnemonic(short)
	require.NoError(t, err)
	require.Equal(t, mnemonic, got)
}

func Test_FromShortToMnemonic(t *testing.T) {
	short := "1692.861.427.1190.1139.287.1329.2040.1220.1244.715.941."
	mnemonic, err := FromShortToMnemonic(short)
	require.NoError(t, err)
	require.Equal(t, "squeeze hill cube network mobile catalog plate yellow obtain oppose floor interest", mnemonic)
}
