package panacea_client

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/go-bip39"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"os"
	"strconv"
	"testing"
)

var (
	coinType = uint32(371)

	testnetMnemonic = os.Getenv("TESTNET_MNEMONIC")

	validatorPrivKey, _ = generatePrivateKeyFromMnemonic(testnetMnemonic, coinType, 0, 0)
)

func generatePrivateKeyFromMnemonic(mnemonic string, coinType, account, index uint32) (secp256k1.PrivKey, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	hdPath := hd.NewFundraiserParams(account, coinType, index).String()
	master, ch := hd.ComputeMastersFromSeed(bip39.NewSeed(mnemonic, ""))

	return hd.DerivePrivateKeyForPath(master, ch, hdPath)
}

func TestMakeCertsFile(t *testing.T) {
	privKey, err := generatePrivateKeyFromMnemonic(testnetMnemonic, 371, 0, 0)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		sellerPriv, err := generatePrivateKeyFromMnemonic(testnetMnemonic, coinType, 1, uint32(i))
		require.NoError(t, err)

		addr, err := bech32.ConvertAndEncode("panacea", sellerPriv.PubKey().Address().Bytes())
		require.NoError(t, err)

		jsonInput := []byte(`{
		"name": "This is a name` + strconv.Itoa(i) + `",
		"description": "This is a description, man",
		"body": [{ "type": "markdown", "attributes": { "value": "` + addr + `" } }]
	}`)
		fmt.Println(string(jsonInput))
	}

	fmt.Println(bech32.ConvertAndEncode("panacea", privKey.PubKey().Address().Bytes()))

}
