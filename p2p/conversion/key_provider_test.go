package conversion

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	. "gopkg.in/check.v1"
)

const testPriKey = "OTI4OTdkYzFjMWFhMjU3MDNiMTE4MDM1OTQyY2Y3MDVkOWFhOGIzN2JlOGIwOWIwMTZjYTkxZjNjOTBhYjhlYQ=="

type KeyProviderTestSuite struct{}

var _ = Suite(&KeyProviderTestSuite{})

func TestGetPubKeysFromPeerIDs(t *testing.T) {
	SetupBech32Prefix()
	input := []string{
		"16Uiu2HAmBdJRswX94UwYj6VLhh4GeUf9X3SjBRgTqFkeEMLmfk2M",
		"16Uiu2HAkyR9dsFqkj1BqKw8ZHAUU2yur6ZLRJxPTiiVYP5uBMeMG",
	}
	result, err := GetPubKeysFromPeerIDs(input)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	assert.Len(t, result, 2)
	assert.Equal(t, "thorpub1addwnpepqtctt9l4fddeh0krvdpxmqsxa5z9xsa0ac6frqfhm9fq6c6u5lck5s8fm4n", result[0])
	assert.Equal(t, "thorpub1addwnpepqga5cupfejfhtw507sh36fvwaekyjt5kwaw0cmgnpku0at2a87qqkp60t43", result[1])
	input = append(input, "whatever")
	result, err = GetPubKeysFromPeerIDs(input)
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func (*KeyProviderTestSuite) TestGetPriKey(c *C) {
	pk, err := GetPriKey("whatever")
	c.Assert(err, NotNil)
	c.Assert(pk, IsNil)
	input := base64.StdEncoding.EncodeToString([]byte("whatever"))
	pk, err = GetPriKey(input)
	c.Assert(err, NotNil)
	c.Assert(pk, IsNil)
	pk, err = GetPriKey(testPriKey)
	c.Assert(err, IsNil)
	c.Assert(pk, NotNil)
	result, err := GetPriKeyRawBytes(pk)
	c.Assert(err, IsNil)
	c.Assert(result, NotNil)
	c.Assert(result, HasLen, 32)
}

//
//func (KeyProviderTestSuite) TestGetPeerIDs(c *C) {
//	pubKeys := []string{
//		"thorpub1addwnpepqtctt9l4fddeh0krvdpxmqsxa5z9xsa0ac6frqfhm9fq6c6u5lck5s8fm4n",
//		"thorpub1addwnpepqga5cupfejfhtw507sh36fvwaekyjt5kwaw0cmgnpku0at2a87qqkp60t43",
//	}
//	peers, err := GetPeerIDs(pubKeys)
//	c.Assert(err, IsNil)
//	c.Assert(peers, NotNil)
//	c.Assert(peers, HasLen, 2)
//	c.Assert(peers[0].String(), Equals, "16Uiu2HAmBdJRswX94UwYj6VLhh4GeUf9X3SjBRgTqFkeEMLmfk2M")
//	c.Assert(peers[1].String(), Equals, "16Uiu2HAkyR9dsFqkj1BqKw8ZHAUU2yur6ZLRJxPTiiVYP5uBMeMG")
//	pubKeys = append(pubKeys, "helloworld")
//	peers, err = GetPeerIDs(pubKeys)
//	c.Assert(err, NotNil)
//	c.Assert(peers, IsNil)
//}
//
//func (KeyProviderTestSuite) TestGetPeerIDFromPubKey(c *C) {
//	pID, err := GetPeerIDFromPubKey("thorpub1addwnpepqtctt9l4fddeh0krvdpxmqsxa5z9xsa0ac6frqfhm9fq6c6u5lck5s8fm4n")
//	c.Assert(err, IsNil)
//	c.Assert(pID.String(), Equals, "16Uiu2HAmBdJRswX94UwYj6VLhh4GeUf9X3SjBRgTqFkeEMLmfk2M")
//	pID1, err := GetPeerIDFromPubKey("whatever")
//	c.Assert(err, NotNil)
//	c.Assert(pID1.String(), Equals, "")
//}

func (KeyProviderTestSuite) TestCheckKeyOnCurve(c *C) {
	_, err := CheckKeyOnCurve("aa")
	c.Assert(err, NotNil)
	SetupBech32Prefix()
	_, err = CheckKeyOnCurve("thorpub1addwnpepqtctt9l4fddeh0krvdpxmqsxa5z9xsa0ac6frqfhm9fq6c6u5lck5s8fm4n")
	c.Assert(err, IsNil)
}

func Test_GetPeerIDFromPubKeyByEth(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "59fb9b933e036204bbe24d3e9d2215952704341159bb6d15526e195b7e250a990dedff123dfa4ae5da1273f5d9863c9389bfc4c53e94b250d565c1669ba3bb17",
		},
		// {
		// 	desc: "0336e12f4a3a175086bbd983f050453194810d8afcfec4e08a898c8e59778c649f",
		// },
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ret, err := GetPeerIDFromPubKeyByEth(tC.desc)
			if ret != "" {
				t.Log("ret --------- ", ret)
			}
			if err != nil {
				t.Log("err --------- ", err)
			}

		})
	}
}

func Test_GetPubKeyFromPeerIDByEth(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "16Uiu2HAmJiFAAggjvwMftkM8nGrAutb6LYqCQYeUFk23WCB2dsiQ",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ret, err := GetPubKeyFromPeerIDByEth(tC.desc)
			if ret != "" {
				t.Log("ret --------- ", ret)
			}
			if err != nil {
				t.Log("err --------- ", err)
			}

		})
	}
}
