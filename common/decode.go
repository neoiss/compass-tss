package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	xrp "github.com/Peersyst/xrpl-go/address-codec"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	dogechaincfg "github.com/eager7/dogd/chaincfg"
	"github.com/eager7/dogutil"
	ethcommon "github.com/ethereum/go-ethereum/common"

	pkgaddress "github.com/mapprotocol/compass-tss/pkg/address"
)

func evmAddressToBytes(address string) ([]byte, error) {
	addr := TrimHexPrefix(address)

	if len(addr) == 0 {
		return []byte{}, fmt.Errorf("empty evm address")
	}

	if !ethcommon.IsHexAddress(addr) {
		return []byte{}, fmt.Errorf("invalid evm address: %s", address)
	}

	decoded, err := hex.DecodeString(addr)
	if err != nil {
		return []byte{}, fmt.Errorf("fail to decode hex string: %w", err)
	}
	return decoded, nil
}

// tronAddressToBytes converts a tron address to bytes
// tron address format: version(1 bytes)+address(20 bytes)+checksum(4 bytes)
// version: 0x41(mainnet)
func tronAddressToBytes(address string) ([]byte, error) {
	if len(address) == 0 {
		return []byte{}, fmt.Errorf("empty tron address")
	}

	decoded := base58.Decode(address)
	if len(decoded) != 25 {
		return []byte{}, fmt.Errorf("invalid tron address length")
	}
	if decoded[0] != 0x41 {
		return []byte{}, fmt.Errorf("invalid tron address")
	}

	// verify checksum
	hash := sha256.Sum256(decoded[:21])
	hash2 := sha256.Sum256(hash[:])
	expectedChecksum := hash2[:4]
	actualChecksum := decoded[21:25]

	if !bytes.Equal(expectedChecksum, actualChecksum) {
		return []byte{}, fmt.Errorf("invalid tron address checksum")
	}

	return decoded[1:21], nil
}

func solanaAddressToBytes(address string) ([]byte, error) {
	decoded := base58.Decode(address)
	if len(decoded) != 32 {
		return []byte{}, fmt.Errorf("invalid solana address length")
	}
	return decoded, nil
}

func xrpAddressToBytes(address string) ([]byte, error) {
	// checks checksum and returns prefix (1 byte, 0x00) + account id (20 bytes)
	decoded, err := xrp.Base58CheckDecode(address)
	if err != nil {
		return nil, fmt.Errorf("invalid xrp address: %w", err)
	}

	if len(decoded) != 21 || decoded[0] != 0x00 {
		return nil, fmt.Errorf("invalid xrp address")
	}

	// return
	return xrp.DecodeBase58(address), nil
}

func bitcoinAddressToBytes(address string) ([]byte, error) {
	chainNetwork := CurrentChainNetwork

	var net *chaincfg.Params
	switch chainNetwork {
	case TestNet:
		net = &chaincfg.TestNet3Params
	case MainNet:
		net = &chaincfg.MainNetParams
	}

	btcAddr, err := btcutil.DecodeAddress(address, net)
	if err != nil {
		return nil, fmt.Errorf("invalid bitcoin address: %w", err)
	}

	decoded, err := pkgaddress.EncodeBitcoinAddressToBytes(btcAddr)
	if err != nil {
		return nil, fmt.Errorf("fail to decode bitcoin hex address: %w", err)
	}
	return decoded, nil
}

func dogeAddressToBytes(address string) ([]byte, error) {
	chainNetwork := CurrentChainNetwork

	var net *dogechaincfg.Params
	switch chainNetwork {
	case TestNet:
		net = &dogechaincfg.TestNet3Params
	case MainNet:
		net = &dogechaincfg.MainNetParams
	}

	btcAddr, err := dogutil.DecodeAddress(address, net)
	if err != nil {
		return nil, fmt.Errorf("invalid bitcoin address: %w", err)
	}

	decoded, err := pkgaddress.EncodeDOGEAddressToBytes(btcAddr)
	if err != nil {
		return nil, fmt.Errorf("fail to decode bitcoin hex address: %w", err)
	}
	return decoded, nil
}
