package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	xrp "github.com/Peersyst/xrpl-go/address-codec"
	"github.com/btcsuite/btcd/btcutil/base58"
	ethcommon "github.com/ethereum/go-ethereum/common"
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
