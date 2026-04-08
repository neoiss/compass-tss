package tron

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/mapprotocol/compass-tss/common"
	keys2 "github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"

	"github.com/ethereum/go-ethereum/crypto"
)

type KeyManager struct {
	priv *ecdsa.PrivateKey
	pub  common.PubKey
}

func NewLocalKeyManager(
	keys *keys2.Keys,
) (*KeyManager, error) {
	priv, err := keys.GetPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("fail to get private key: %w", err)
	}

	evmPrivateKey, err := evm.GetPrivateKey(priv)
	if err != nil {
		return nil, err
	}
	compressPkBytes := crypto.CompressPubkey(&evmPrivateKey.PublicKey)
	pk, err := common.NewPubKey(hex.EncodeToString(compressPkBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to get pub key: %w", err)
	}

	return &KeyManager{
		priv: evmPrivateKey,
		pub:  pk,
	}, nil
}

func (m *KeyManager) Pubkey() common.PubKey {
	return m.pub
}

func (m *KeyManager) Sign(msg []byte) ([]byte, error) {
	return crypto.Sign(msg, m.priv)
}
