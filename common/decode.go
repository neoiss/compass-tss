package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"reflect"
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
	if len(decoded) != 20 {
		return []byte{}, fmt.Errorf("invalid evm address length")
	}
	return decoded, nil
}

func tronAddressToBytes(address string) ([]byte, error) {
	if len(address) == 0 {
		return []byte{}, fmt.Errorf("empty tron address")
	}

	decoded := base58.Decode(address)
	if len(decoded) != 25 {
		return []byte{}, fmt.Errorf("invalid tron address length")
	}
	if decoded[1] != 0x41 {
		return []byte{}, fmt.Errorf("invalid tron address")
	}

	// verify checksum
	hash := sha256.Sum256(decoded[:21])
	hash2 := sha256.Sum256(hash[:])
	expectedChecksum := hash2[:4]
	actualChecksum := decoded[21:25]

	if !reflect.DeepEqual(expectedChecksum, actualChecksum) {
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
