package address

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	P2WPKH = "0x00"
	P2WSH  = "0x00"
	P2PKH  = "0x01"
	P2TR   = "0x01"
	P2SH   = "0x05"
)

const (
	P2WPKHOrP2WSH = "0x00"
	P2PKHOrP2TR   = "0x01"
)

var (
	P2WPKHBytes = []byte{0x00}
	P2WSHBytes  = []byte{0x00}
	P2PKHBytes  = []byte{0x01}
	P2TRBytes   = []byte{0x01}
	P2SHBytes   = []byte{0x05}
)

var ErrUnknownAddressType = errors.New("unknown address type")

type unsupportedPublicKeyError struct {
	prefix string
	length int
}

func newUnsupportedPublicKeyLenError(prefix string, length int) unsupportedPublicKeyError {
	return unsupportedPublicKeyError{
		prefix: prefix,
		length: length,
	}
}

func (e unsupportedPublicKeyError) Error() string {
	return fmt.Sprintf("unsupported public key, prefix: %s, length: %d", e.prefix, e.length)
}

type addrPrefix struct {
	str   string
	bytes []byte
}

func getBTCAddrPrefix(addr btcutil.Address) (addrPrefix, bool) {
	switch addr.(type) {
	case *btcutil.AddressWitnessPubKeyHash:
		return addrPrefix{str: P2WPKH, bytes: P2WPKHBytes}, true
	case *btcutil.AddressTaproot:
		return addrPrefix{str: P2TR, bytes: P2TRBytes}, true
	case *btcutil.AddressWitnessScriptHash:
		return addrPrefix{str: P2WSH, bytes: P2WSHBytes}, true
	case *btcutil.AddressPubKeyHash:
		return addrPrefix{str: P2PKH, bytes: P2PKHBytes}, true
	case *btcutil.AddressScriptHash:
		return addrPrefix{str: P2SH, bytes: P2SHBytes}, true
	default:
		return addrPrefix{}, false
	}
}

func EncodeBitcoinAddress(addr btcutil.Address) (string, error) {
	if addr == nil {
		return "", fmt.Errorf("address cannot be nil")
	}

	prefix, ok := getBTCAddrPrefix(addr)
	if !ok {
		return "", ErrUnknownAddressType
	}
	return prefix.str + hex.EncodeToString(addr.ScriptAddress()), nil
}

func EncodeBitcoinAddressToBytes(addr btcutil.Address) ([]byte, error) {
	if addr == nil {
		return nil, fmt.Errorf("address cannot be nil")
	}

	prefix, ok := getBTCAddrPrefix(addr)
	if !ok {
		return nil, ErrUnknownAddressType
	}
	script := addr.ScriptAddress()
	result := make([]byte, 0, len(prefix.bytes)+len(script))
	result = append(result, prefix.bytes...)
	result = append(result, script...)
	return result, nil
}

func DecodeBitcoinAddress(addr string, network *chaincfg.Params) (btcutil.Address, error) {
	if !HasHexPrefix(addr) {
		addr = "0x" + addr
	}
	if len(addr) <= 4 {
		return nil, fmt.Errorf("invalid address: %s", addr)
	}

	prefix := addr[:4]
	publicKey := addr[4:]
	publicKeyLen := len(publicKey)
	publicKeyBytes := ethcommon.Hex2Bytes(publicKey)

	switch prefix {
	case P2WPKHOrP2WSH:
		switch publicKeyLen {
		case 40: // P2WPKH
			return btcutil.NewAddressWitnessPubKeyHash(publicKeyBytes, network)
		case 64: // P2WSH
			return btcutil.NewAddressWitnessScriptHash(publicKeyBytes, network)
		default:
			return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
		}
	case P2PKHOrP2TR:
		switch publicKeyLen {
		case 40: // P2PKH
			return btcutil.NewAddressPubKeyHash(publicKeyBytes, network) // base 58
		case 64: // P2TR
			return btcutil.NewAddressTaproot(publicKeyBytes, network)
		default:
			return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
		}
	case P2SH:
		return btcutil.NewAddressScriptHashFromHash(publicKeyBytes, network) // base 58
	default:
		return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
	}
}

func HasHexPrefix(s string) bool {
	return len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X")
}
