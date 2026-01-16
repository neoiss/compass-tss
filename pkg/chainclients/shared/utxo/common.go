package utxo

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/mapprotocol/compass-tss/common"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
)

var (
	bytesType, _   = abi.NewType("bytes", "bytes", nil)
	addressType, _ = abi.NewType("address", "address", nil)
	uint256Type, _ = abi.NewType("uint256", "uint256", nil)
)

func GetAsgardAddress(chain common.Chain, bridge shareTypes.Bridge) ([]common.Address, error) {
	vaults, err := bridge.GetAsgardPubKeys()
	if err != nil {
		return nil, fmt.Errorf("fail to get asgards : %w", err)
	}

	newAddresses := make([]common.Address, 0)
	for _, v := range vaults {
		var addr common.Address
		addr, err = v.CompressedPubKey.GetAddress(chain)
		if err != nil {
			continue
		}
		newAddresses = append(newAddresses, addr)
	}
	return newAddresses, nil
}

func GetAsgardPubKeyByAddress(chain common.Chain, bridge shareTypes.Bridge, address common.Address) ([]byte, error) {
	vaults, err := bridge.GetAsgardPubKeys()
	if err != nil {
		return nil, fmt.Errorf("fail to get asgards : %w", err)
	}

	for _, v := range vaults {
		addr, err := v.CompressedPubKey.GetAddress(chain)
		if err != nil {
			continue
		}
		if !address.Equals(addr) {
			continue
		}

		pubKey, err := hex.DecodeString(strings.TrimPrefix(v.PubKey.String(), "04"))
		if err != nil {
			return nil, fmt.Errorf("fail to decode pubkey(%s)", v.PubKey.String())
		}
		return pubKey, nil
	}
	return nil, fmt.Errorf("fail to get asgard pub key by address(%s)", address)
}

//func GetAsgardAddress2PubKeyMapped(chain common.Chain, bridge shareTypes.Bridge) (map[common.Address][]byte, error) {
//	vaults, err := bridge.GetAsgardPubKeys()
//	if err != nil {
//		return nil, fmt.Errorf("fail to get asgards : %w", err)
//	}
//
//	addr2pub := make(map[common.Address][]byte, 0)
//	for _, v := range vaults {
//		addr, err := v.CompressedPubKey.GetAddress(chain)
//		if err != nil {
//			continue
//		}
//		pubKey, err := hex.DecodeString(strings.TrimPrefix(v.PubKey.String(), "04"))
//		if err != nil {
//			continue
//		}
//		addr2pub[addr] = pubKey
//	}
//	return addr2pub, nil
//}

type Affiliate struct {
	ID  uint16
	Bps uint16
}

func EncodeAffiliateData(affiliates []*Affiliate) ([]byte, error) {
	if len(affiliates) == 0 {
		return []byte{}, nil
	}

	buf := make([]byte, len(affiliates)*4)
	for i, affiliate := range affiliates {
		if affiliate == nil {
			return nil, fmt.Errorf("affiliate is nil")
		}
		offset := i * 4

		binary.BigEndian.PutUint16(buf[offset:], affiliate.ID)
		binary.BigEndian.PutUint16(buf[offset+2:], affiliate.Bps)
	}
	return buf, nil
}

func EncodeRelayData(token ethcommon.Address, minAmount *big.Int) ([]byte, error) {
	args := abi.Arguments{
		{Type: addressType},
		{Type: uint256Type},
	}
	packed, err := args.Pack(token, minAmount)
	if err != nil {
		return nil, err
	}
	return packed, nil
}

func EncodeTargetData(to []byte, toChain *big.Int) ([]byte, error) {
	args := abi.Arguments{
		{Type: bytesType},
		{Type: uint256Type},
	}
	packed, err := args.Pack(to, toChain)
	if err != nil {
		return nil, err
	}
	return packed, nil
}

func EncodePayload(affiliateData, relayData, targetData []byte) ([]byte, error) {
	args := abi.Arguments{
		{Type: bytesType},
		{Type: bytesType},
		{Type: bytesType},
	}
	packed, err := args.Pack(affiliateData, relayData, targetData)
	if err != nil {
		return nil, err
	}
	return packed, nil
}

func ConvertDecimal(amount *big.Int, srcDecimal uint64, dstDecimal uint64) *big.Int {
	dstAmount := amount
	if srcDecimal > dstDecimal {
		exp := new(big.Int).Exp(big.NewInt(10), new(big.Int).SetUint64(srcDecimal-dstDecimal), nil)
		dstAmount = new(big.Int).Div(amount, exp)
	} else if srcDecimal < dstDecimal {
		exp := new(big.Int).Exp(big.NewInt(10), new(big.Int).SetUint64(dstDecimal-srcDecimal), nil)
		dstAmount = new(big.Int).Mul(amount, exp)
	}
	return dstAmount
}
