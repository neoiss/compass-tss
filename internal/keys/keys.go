package keys

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"io"
	"os/user"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// folder name for thorchain thorcli
	thorchainCliFolderName = `.compass`
)

// Keys manages all the keys used by thorchain
type Keys struct {
	signerName string
	password   string // TODO this is a bad way , need to fix it
	kb         ckeys.Keyring
}

// NewKeysWithKeybase create a new instance of Keys
func NewKeysWithKeybase(kb ckeys.Keyring, name, password string) *Keys {
	return &Keys{
		signerName: name,
		password:   password,
		kb:         kb,
	}
}

// GetKeyringKeybase return keyring and key info
func GetKeyringKeybase(priKeyStr, signerName string) (ckeys.Keyring, *ckeys.Record, error) {
	//if len(signerName) == 0 {
	//	return nil, nil, fmt.Errorf("signer name is empty")
	//}
	//if len(password) == 0 {
	//	return nil, nil, fmt.Errorf("password is empty")
	//}
	if len(priKeyStr) == 0 {
		return nil, nil, fmt.Errorf("priKey is empty")
	}
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	kb := ckeys.NewInMemory(cdc)
	err := kb.ImportPrivKeyHex(signerName, priKeyStr, string(hd.Secp256k1.Name()))
	if err != nil {
		return nil, nil, err
	}
	//
	//buf := bytes.NewBufferString(password)
	//// the library used by keyring is using ReadLine , which expect a new line
	//buf.WriteByte('\n')
	//kb, err := getKeybase(chainHomeFolder, buf)
	//if err != nil {
	//	return nil, nil, fmt.Errorf("fail to get keybase,err:%w", err)
	//}
	//// the keyring library which used by cosmos sdk , will use interactive terminal if it detect it has one
	//// this will temporary trick it think there is no interactive terminal, thus will read the password from the buffer provided
	//oldStdIn := os.Stdin
	//defer func() {
	//	os.Stdin = oldStdIn
	//}()
	//os.Stdin = nil
	// keyring
	// name -> key  private-> value
	//si, err := kb.Key(signerName) // first-cosmos ｜ second-gen
	//if err != nil {
	//	return nil, nil, fmt.Errorf("fail to get signer info(%s): %w", signerName, err)
	//}
	return kb, nil, nil
}

// getKeybase will create an instance of Keybase
func getKeybase(thorchainHome string, reader io.Reader) (ckeys.Keyring, error) {
	cliDir := thorchainHome
	if len(thorchainHome) == 0 {
		usr, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("fail to get current user,err:%w", err)
		}
		cliDir = filepath.Join(usr.HomeDir, thorchainCliFolderName)
	}

	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	return ckeys.New(sdk.KeyringServiceName(), ckeys.BackendFile, cliDir, reader, cdc)
}

// GetSignerInfo return signer info
func (k *Keys) GetSignerInfo() *ckeys.Record {
	record, err := k.kb.Key(k.signerName)
	if err != nil {
		panic(err)
	}
	return record
}

// GetPrivateKey return the private key
func (k *Keys) GetPrivateKey() (cryptotypes.PrivKey, error) {
	//addr, err := sdk.AccAddressFromHexUnsafe("2b7588165556aB2fA1d30c520491C385BAa424d8")
	//if err != nil {
	//	return nil, err
	//}
	//privKeyArmor, err := k.kb.ExportPrivKeyArmorByAddress(addr, k.password)
	privKeyArmor, err := k.kb.ExportPrivKeyArmor(k.signerName, k.password)
	if err != nil {
		return nil, err
	}
	priKey, _, err := crypto.UnarmorDecryptPrivKey(privKeyArmor, k.password)
	if err != nil {
		return nil, fmt.Errorf("fail to unarmor private key: %w", err)
	}
	return priKey, nil
}

// GetKeybase return the keybase
func (k *Keys) GetKeybase() ckeys.Keyring {
	return k.kb
}
