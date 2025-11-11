package common

import (
	"encoding/json"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func HasHexPrefix(s string) bool {
	return len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X")
}

func TrimHexPrefix(s string) string {
	if len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X") {
		return s[2:]
	}
	return s
}

func JSON(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}

func CompressPubKey(pks []byte) (string, error) {
	if len(pks) != 64 {
		return "", fmt.Errorf("invalid pub key, length(%d)", len(pks))
	}
	pk, err := ethcrypto.UnmarshalPubkey(append([]byte{4}, pks...))
	if err != nil {
		return "", err
	}
	cpkBytes := ethcrypto.CompressPubkey(pk)
	return ethcommon.Bytes2Hex(cpkBytes), nil
}
