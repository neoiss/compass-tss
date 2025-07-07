package tss

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/cmd"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

func TestTSSKeyGen(t *testing.T) { TestingT(t) }

type KeyGenTestSuite struct{}

var _ = Suite(&KeyGenTestSuite{})

func (*KeyGenTestSuite) SetUpSuite(c *C) {
}

const (
	signerNameForTest     = `jack`
	signerPasswordForTest = `password`
)

func (*KeyGenTestSuite) setupKeysForTest(c *C) string {
	ns := strconv.Itoa(time.Now().Nanosecond())
	thorcliDir := filepath.Join(os.TempDir(), ns, ".thorcli")
	c.Logf("thorcliDir:%s", thorcliDir)
	buf := bytes.NewBufferString(signerPasswordForTest)
	// the library used by keyring is using ReadLine , which expect a new line
	buf.WriteByte('\n')
	buf.WriteString(signerPasswordForTest)
	buf.WriteByte('\n')
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	kb, err := cKeys.New(cosmos.KeyringServiceName(), cKeys.BackendFile, thorcliDir, buf, cdc)
	c.Assert(err, IsNil)
	info, _, err := kb.NewMnemonic(signerNameForTest, cKeys.English, cmd.THORChainHDPath, signerPasswordForTest, hd.Secp256k1)
	c.Assert(err, IsNil)
	c.Logf("name:%s", info.Name)
	return thorcliDir
}

//func (kts *KeyGenTestSuite) TestNewTssKenGen(c *C) {
//	oldStdIn := os.Stdin
//	defer func() {
//		os.Stdin = oldStdIn
//	}()
//	os.Stdin = nil
//	folder := kts.setupKeysForTest(c)
//	defer func() {
//		err := os.RemoveAll(folder)
//		c.Assert(err, IsNil)
//	}()
//	kb, _, err := mapclient.GetKeyringKeybase(folder, signerNameForTest, signerPasswordForTest)
//	c.Assert(err, IsNil)
//	k := mapclient.NewKeysWithKeybase(kb, signerNameForTest, signerPasswordForTest)
//	c.Assert(k, NotNil)
//	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
//		c.Logf("requestUri:%s", req.RequestURI)
//	}))
//	b, err := mapclient.NewThorchainBridge(config.BifrostClientConfiguration{
//		ChainID:      "thorchain",
//		ChainHost:    server.Listener.Addr().String(),
//		SignerName:   "bob",
//		SignerPasswd: "password",
//	}, nil, k)
//	c.Assert(err, IsNil)
//	kg, err := NewTssKeyGen(k, nil, b)
//	c.Assert(err, IsNil)
//	c.Assert(kg, NotNil)
//}
