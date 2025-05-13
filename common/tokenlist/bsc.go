package tokenlist

import (
	"encoding/json"

	"github.com/mapprotocol/compass-tss/common/tokenlist/bsctokens"
)

var bscTokenList EVMTokenList

func init() {
	if err := json.Unmarshal(bsctokens.BSCTokenListRaw, &bscTokenList); err != nil {
		panic(err)
	}
}

func GetBSCTokenList() EVMTokenList {
	return bscTokenList
}
