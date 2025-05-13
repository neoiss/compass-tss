package tokenlist

import (
	"encoding/json"

	"github.com/mapprotocol/compass-tss/common/tokenlist/basetokens"
)

var baseTokenList EVMTokenList

func init() {
	if err := json.Unmarshal(basetokens.BASETokenListRaw, &baseTokenList); err != nil {
		panic(err)
	}
}

func GetBASETokenList() EVMTokenList {
	return baseTokenList
}
