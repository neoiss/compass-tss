package thorchain

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"testing"
)

const nativeToken = "0x0000000000000000000000000000000000000000"

func TestOutboundMemo(t *testing.T) {
	memo := NewOutboundMemo("97", nativeToken, "0x0Eb16A9cFDf8e3A4471EF190eE63de5A24f38787", "1000", "").String()
	t.Log("memo: ", memo)

	parsedMemo, err := ParseMemo(memo)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("parsed memo: ", parsedMemo.String())

	if memo != parsedMemo.String() {
		t.Errorf("memo != parsedMemo")
	}

	t.Log(ethcommon.Hex2Bytes("0000000000000000000000000000000000000000"))
}
