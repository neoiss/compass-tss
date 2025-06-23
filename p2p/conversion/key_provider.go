package conversion

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	tcrypto "github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types/bech32/legacybech32" // nolint:staticcheck
	ecommon "github.com/ethereum/go-ethereum/common"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

// GetPeerIDFromPubKeyByEth get the peer.ID from public key format node pub key
func GetPeerIDFromPubKeyByEth(pubkey string) (peer.ID, error) {
	ppk, err := crypto.UnmarshalSecp256k1PublicKey(ecommon.Hex2Bytes(pubkey))
	if err != nil {
		return "", fmt.Errorf("fail to convert pubkey to the crypto pubkey used in libp2p: %w", err)
	}
	return peer.IDFromPublicKey(ppk)
}

// GetPubKeyFromPeerIDByEth extract the pub key from PeerID
func GetPubKeyFromPeerIDByEth(pID string) (string, error) {
	peerID, err := peer.Decode(pID)
	if err != nil {
		return "", fmt.Errorf("fail to decode peer id: %w", err)
	}

	pubKey, err := peerID.ExtractPublicKey()
	if err != nil {
		return "", fmt.Errorf("failed to extract public key from peer ID: %w", err)
	}

	rawPubKey, err := pubKey.Raw()
	if err != nil {
		return "", fmt.Errorf("failed to get raw public key: %w", err)
	}

	fmt.Println("rawPubKey ---------------- ", rawPubKey)
	ethPubKey, err := ecrypto.DecompressPubkey(rawPubKey)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal ECDSA public key: %w", err)
	}
	pubBytes := ecrypto.FromECDSAPub(ethPubKey)

	return hex.EncodeToString(pubBytes), nil
}

// GetPeerIDsFromPubKeys convert a list of node pub key to their peer.ID
func GetPeerIDsFromPubKeys(pubkeys []string) ([]peer.ID, error) {
	var peerIDs []peer.ID
	for _, item := range pubkeys {
		peerID, err := GetPeerIDFromPubKeyByEth(item)
		if err != nil {
			return nil, err
		}
		peerIDs = append(peerIDs, peerID)
	}
	return peerIDs, nil
}

// GetPubKeysFromPeerIDs given a list of peer ids, and get a list og pub keys.
func GetPubKeysFromPeerIDs(peers []string) ([]string, error) {
	var result []string
	for _, item := range peers {
		pKey, err := GetPubKeyFromPeerIDByEth(item)
		if err != nil {
			return nil, fmt.Errorf("fail to get pubkey from peerID: %w", err)
		}
		result = append(result, pKey)
	}
	return result, nil
}

func GetPriKey(priKeyString string) (tcrypto.PrivKey, error) {
	priHexBytes, err := base64.StdEncoding.DecodeString(priKeyString)
	if err != nil {
		return nil, fmt.Errorf("fail to decode private key: %w", err)
	}
	rawBytes, err := hex.DecodeString(string(priHexBytes))
	if err != nil {
		return nil, fmt.Errorf("fail to hex decode private key: %w", err)
	}
	var priKey secp256k1.PrivKey = rawBytes[:32]
	return priKey, nil
}

func GetPriKeyRawBytes(priKey tcrypto.PrivKey) ([]byte, error) {
	var keyBytesArray [32]byte
	pk, ok := priKey.(secp256k1.PrivKey)
	if !ok {
		return nil, errors.New("private key is not secp256p1.PrivKey")
	}
	copy(keyBytesArray[:], pk[:])
	return keyBytesArray[:], nil
}

func CheckKeyOnCurve(pk string) (bool, error) {
	pubKey, err := sdk.UnmarshalPubKey(sdk.AccPK, pk) // nolint:staticcheck
	if err != nil {
		return false, fmt.Errorf("fail to parse pub key(%s): %w", pk, err)
	}
	bPk, err := btcec.ParsePubKey(pubKey.Bytes(), btcec.S256())
	if err != nil {
		return false, err
	}
	return isOnCurve(bPk.X, bPk.Y), nil
}
