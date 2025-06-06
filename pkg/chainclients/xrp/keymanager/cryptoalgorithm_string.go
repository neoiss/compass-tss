// Code generated by "stringer -type=CryptoAlgorithm"; DO NOT EDIT.

package keymanager

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SECP256K1-0]
	_ = x[ED25519-1]
}

const _CryptoAlgorithm_name = "SECP256K1ED25519"

var _CryptoAlgorithm_index = [...]uint8{0, 9, 16}

func (i CryptoAlgorithm) String() string {
	if i < 0 || i >= CryptoAlgorithm(len(_CryptoAlgorithm_index)-1) {
		return "CryptoAlgorithm(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CryptoAlgorithm_name[_CryptoAlgorithm_index[i]:_CryptoAlgorithm_index[i+1]]
}
