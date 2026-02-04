package keys

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	ekeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ecommon "github.com/ethereum/go-ethereum/common"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/pkg/keystore"
)

const (
	// folder name for relay thorcli
	relayCliFolderName = `.compass`
)

// Keys manages all the keys used by relay
type Keys struct {
	signerName string
	password   string // TODO this is a bad way , need to fix it
	kb         ckeys.Keyring
	keyStore   *ekeystore.Key
}

// NewKeysWithKeybase create a new instance of Keys
func NewKeysWithKeybase(kb ckeys.Keyring, name, password string, keyStore *ekeystore.Key) *Keys {
	return &Keys{
		signerName: name,
		password:   password,
		kb:         kb,
		keyStore:   keyStore,
	}
}

// GetKeyringKeybase return keyring and key info
func GetKeyringKeybase(keyStorePath, signerName string) (ckeys.Keyring, *ekeystore.Key, error) {
	if len(signerName) == 0 {
		return nil, nil, fmt.Errorf("signer name is empty")
	}
	if len(keyStorePath) == 0 {
		return nil, nil, fmt.Errorf("keyStorePath is empty")
	}

	kpI, err := keystore.DecryptKey(keyStorePath)
	if err != nil {
		return nil, nil, err
	}
	priBytes := ecrypto.FromECDSA(kpI.PrivateKey)
	priKeyStr := common.Bytes2Hex(priBytes)

	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	kb := ckeys.NewInMemory(cdc)
	err = kb.ImportPrivKeyHex(signerName, priKeyStr, string(hd.Secp256k1.Name()))
	if err != nil {
		return nil, nil, err
	}
	return kb, kpI, nil
}

// GetSignerInfo return signer info
func (k *Keys) GetSignerInfo() *ckeys.Record {
	record, err := k.kb.Key(k.signerName)
	if err != nil {
		panic(err)
	}
	return record
}

func (k *Keys) GetEthAddress() (common.Address, error) {
	addr, err := AddressFromPublicKey(&k.keyStore.PrivateKey.PublicKey)
	if err != nil {
		return constants.ZeroAddress, err
	}
	return addr, nil
}

// GetPrivateKey return the private key
func (k *Keys) GetPrivateKey() (cryptotypes.PrivKey, error) {
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

func PrivateKeyFromHex(hexKey string) (*ecdsa.PrivateKey, error) {
	cleaned := strings.TrimSpace(hexKey)
	//cleaned = strings.TrimPrefix(cleaned, "0x")

	if len(cleaned) != 64 {
		return nil, errors.New("invalid hex key")
	}

	privateKeyBytes, err := hex.DecodeString(cleaned)
	if err != nil {
		return nil, fmt.Errorf("decodeString failed: %v", err)
	}

	privateKey, err := ecrypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("ToECDSA failed: %v", err)
	}

	return privateKey, nil
}

func PublicKeyFromPrivate(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	if privateKey == nil {
		return nil, errors.New("invalid private key")
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key")
	}

	return publicKeyECDSA, nil
}

func AddressFromPublicKey(publicKey *ecdsa.PublicKey) (common.Address, error) {
	if publicKey == nil {
		return common.Address{}, errors.New("invalid public key")
	}

	publicKeyBytes := ecrypto.FromECDSAPub(publicKey)
	hash := ecrypto.Keccak256(publicKeyBytes[1:])
	address := common.BytesToAddress(hash[12:])

	return address, nil
}

func GetAddressByCompressPk(compressPk string) (common.Address, error) {
	if compressPk == "" {
		return common.Address{}, errors.New("empty public key")
	}
	pk, err := ecrypto.DecompressPubkey(ecommon.Hex2Bytes(compressPk))
	if err != nil {
		return common.Address{}, err
	}
	return AddressFromPublicKey(pk)
}
