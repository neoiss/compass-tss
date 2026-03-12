package address

import (
	"encoding/hex"
	"fmt"
	"github.com/eager7/dogd/chaincfg"
	"github.com/eager7/dogutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func getDOGEAddrPrefix(addr dogutil.Address) (addrPrefix, bool) {
	switch addr.(type) {
	case *dogutil.AddressWitnessPubKeyHash:
		return addrPrefix{str: P2WPKH, bytes: P2WPKHBytes}, true
	//case *dogutil.AddressTaproot:
	case *dogutil.AddressWitnessScriptHash:
		return addrPrefix{str: P2WSH, bytes: P2WSHBytes}, true
	case *dogutil.AddressPubKeyHash:
		return addrPrefix{str: P2PKH, bytes: P2PKHBytes}, true
	case *dogutil.AddressScriptHash:
		return addrPrefix{str: P2SH, bytes: P2SHBytes}, true
	default:
		return addrPrefix{}, false
	}
}

func EncodeDOGEAddress(addr dogutil.Address) (string, error) {
	if addr == nil {
		return "", fmt.Errorf("address cannot be nil")
	}

	prefix, ok := getDOGEAddrPrefix(addr)
	if !ok {
		return "", ErrUnknownAddressType
	}
	return prefix.str + hex.EncodeToString(addr.ScriptAddress()), nil
}

func EncodeDOGEAddressToBytes(addr dogutil.Address) ([]byte, error) {
	if addr == nil {
		return nil, fmt.Errorf("address cannot be nil")
	}

	prefix, ok := getDOGEAddrPrefix(addr)
	if !ok {
		return nil, ErrUnknownAddressType
	}
	script := addr.ScriptAddress()
	result := make([]byte, 0, len(prefix.bytes)+len(script))
	result = append(result, prefix.bytes...)
	result = append(result, script...)
	return result, nil
}

func DecodeDOGEAddress(addr string, network *chaincfg.Params) (dogutil.Address, error) {
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
			return dogutil.NewAddressWitnessPubKeyHash(publicKeyBytes, network)
		case 64: // P2WSH
			return dogutil.NewAddressWitnessScriptHash(publicKeyBytes, network)
		default:
			return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
		}
	case P2PKHOrP2TR:
		switch publicKeyLen {
		case 40: // P2PKH
			return dogutil.NewAddressPubKeyHash(publicKeyBytes, network) // base 58
		//case 64: // P2TR dogecoin not support taproot
		default:
			return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
		}
	case P2SH:
		return dogutil.NewAddressScriptHashFromHash(publicKeyBytes, network) // base 58
	default:
		return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
	}
}
