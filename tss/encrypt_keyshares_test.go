package tss

import (
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"strings"

	"github.com/itchio/lzma"
	"github.com/mapprotocol/compass-tss/p2p/storage"
	. "gopkg.in/check.v1"
)

// -------------------------------------------------------------------------------------
// Setup
// -------------------------------------------------------------------------------------

type EncryptKeySharesSuite struct{}

var _ = Suite(&EncryptKeySharesSuite{})

const (
	LocalStateTestFile = "localstate-test.json"
	Mnemonic           = "profit used piece repeat real curtain endorse tennis tenant sentence include glass return learn upgrade apple crane polar attend before ripple doctor decrease depend"
)

// -------------------------------------------------------------------------------------
// Tests
// -------------------------------------------------------------------------------------

func (s *EncryptKeySharesSuite) TestEncryptKeySharesEmptyPassphrase(c *C) {
	ks, err := EncryptKeyShares("", "")
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "failed keyshare encrypt: signer seed phrase is not set")
	c.Assert(ks, IsNil)
}

func (s *EncryptKeySharesSuite) TestEncryptKeySharesBadMnemonic(c *C) {
	ks, err := EncryptKeyShares("", Mnemonic+" dog")
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "failed keyshare encrypt: signer seed phrase is not 24 words")
	c.Assert(ks, IsNil)

	ks, err = EncryptKeyShares("", "a b c d e f g h i j k l m n o p q r s t u v w x")
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "failed keyshare encrypt: signer seed phrase is not valid bip39 mnemonic")
	c.Assert(ks, IsNil)
}

func (s *EncryptKeySharesSuite) TestEncryptKeySharesBadMnemonicEntropy(c *C) {
	ks, err := EncryptKeyShares("", "dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog")
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "failed keyshare encrypt: signer seed phrase failed entropy check")
	c.Assert(ks, IsNil)
}

func (s *EncryptKeySharesSuite) TestEncryptKeySharesMissingFile(c *C) {
	ks, err := EncryptKeyShares("foo.json", Mnemonic)
	c.Assert(err, NotNil)
	c.Assert(strings.HasPrefix(err.Error(), "failed keyshare encrypt - cannot open key file"), Equals, true)
	c.Assert(ks, IsNil)
}

func (s *EncryptKeySharesSuite) TestEncryptKeySharesCompression(c *C) {
	// encrypt
	ks, err := EncryptKeyShares(LocalStateTestFile, Mnemonic)
	c.Assert(err, IsNil)
	c.Assert(ks, NotNil)

	// ensure we achieved expected compression ratio
	fi, err := os.Stat(LocalStateTestFile)
	c.Assert(err, IsNil)
	if float64(len(ks))/float64(fi.Size()) > 0.4 {
		c.Fatalf("compression ratio over expected: %f", float64(len(ks))/float64(fi.Size()))
	}
}

func (s *EncryptKeySharesSuite) TestEncryptKeyShares(c *C) {
	// encrypt
	ks, err := EncryptKeyShares(LocalStateTestFile, Mnemonic)
	c.Assert(err, IsNil)
	c.Assert(ks, NotNil)

	// decrypt with bad passphrase should fail
	dec, err := DecryptKeyShares(ks, Mnemonic+" y")
	c.Assert(err, NotNil)
	c.Assert(dec, IsNil)

	// decrypt with good passphrase should succeed
	dec, err = DecryptKeyShares(ks, Mnemonic)
	c.Assert(err, IsNil)
	cmpOut := bytes.NewBuffer(dec)
	out := lzma.NewReader(cmpOut)

	// decrypted value should match
	var original, decrypted storage.KeygenLocalState
	f, err := os.Open(LocalStateTestFile)
	c.Assert(err, IsNil)
	defer f.Close()
	err = json.NewDecoder(f).Decode(&original)
	c.Assert(err, IsNil)
	err = json.NewDecoder(out).Decode(&decrypted)
	c.Assert(err, IsNil)
	c.Assert(decrypted, DeepEquals, original)
}

func (s *EncryptKeySharesSuite) TestEncryptKeyShares2(c *C) {
	// encrypt
	ks, err := EncryptKeyShares("localstate-test-2.json", Mnemonic)
	c.Assert(err, IsNil)
	c.Assert(ks, NotNil)
	c.Log("encrypted key shares: ", common.Bytes2Hex(ks))
}

func (s *EncryptKeySharesSuite) TestSaltAndHash(c *C) {
	hash := saltAndHash("foo", 1)
	c.Assert(len(hash), Equals, 32)

	hash2 := saltAndHash("foo", 2)
	c.Assert(len(hash2), Equals, 32)
	c.Assert(hash, Not(Equals), hash2)
}
