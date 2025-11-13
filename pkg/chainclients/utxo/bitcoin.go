package utxo

import (
	"encoding/hex"
	"errors"
	"fmt"

	btcjson "github.com/btcsuite/btcd/btcjson"
	btcchaincfg "github.com/btcsuite/btcd/chaincfg"
	btcwire "github.com/btcsuite/btcd/wire"
	btctxscript "github.com/mapprotocol/compass-tss/txscript/txscript"

	"github.com/mapprotocol/compass-tss/common"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	P2WPKH = "0x00"
	P2WSH  = "0x00"
	P2PKH  = "0x01"
	P2TR   = "0x01"
	P2SH   = "0x05"
)

const (
	P2WPKHOrP2WSH = "0x00"
	P2PKHOrP2TR   = "0x01"
)

var errUnknownAddressType = errors.New("unknown address type")

type unsupportedPublicKeyError struct {
	prefix string
	length int
}

func newUnsupportedPublicKeyLenError(prefix string, length int) unsupportedPublicKeyError {
	return unsupportedPublicKeyError{
		prefix: prefix,
		length: length,
	}
}

func (e unsupportedPublicKeyError) Error() string {
	return fmt.Sprintf("unsupported public key, prefix: %s, length: %d", e.prefix, e.length)
}

func EncodeBitcoinAddress(addr btcutil.Address) (string, error) {
	script := hex.EncodeToString(addr.ScriptAddress())

	switch addr.(type) {
	case *btcutil.AddressWitnessPubKeyHash: // P2WPKH
		return P2WPKH + script, nil
	case *btcutil.AddressTaproot: // P2TR
		return P2TR + script, nil
	case *btcutil.AddressWitnessScriptHash: // P2WSH
		return P2WSH + script, nil
	case *btcutil.AddressPubKeyHash: // P2PKH
		return P2PKH + script, nil
	case *btcutil.AddressScriptHash: // P2SH
		return P2SH + script, nil
	default:
		return "", errUnknownAddressType
	}
}

func DecodeBitcoinAddress(addr string, network *chaincfg.Params) (btcutil.Address, error) {
	if !common.HasHexPrefix(addr) {
		addr = "0x" + addr
	}

	prefix := addr[:4]
	publicKey := addr[4:]
	publicKeyLen := len(publicKey)
	publicKeyBytes := ethcommon.Hex2Bytes(publicKey)

	switch prefix {
	case P2WPKHOrP2WSH:
		switch publicKeyLen {
		case 40: // P2WPKH
			return btcutil.NewAddressWitnessPubKeyHash(publicKeyBytes, network)
		case 64: // P2WSH
			return btcutil.NewAddressWitnessScriptHash(publicKeyBytes, network)
		default:
			return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
		}
	case P2PKHOrP2TR:
		switch publicKeyLen {
		case 40: // P2PKH
			return btcutil.NewAddressPubKeyHash(publicKeyBytes, network) // base 58
		case 64: // P2TR
			return btcutil.NewAddressTaproot(publicKeyBytes, network)
		default:
			return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
		}

	case P2SH:
		return btcutil.NewAddressScriptHashFromHash(publicKeyBytes, network) // base 58
	default:
		return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
	}

}

func (c *Client) getChainCfgBTC() *btcchaincfg.Params {
	switch common.CurrentChainNetwork {
	case common.TestNet:
		return &btcchaincfg.TestNet3Params
	case common.MainNet:
		return &btcchaincfg.MainNetParams
	default:
		c.log.Fatal().Msg("unsupported network")
		return nil
	}
}

func (c *Client) signUTXOBTC(redeemTx *btcwire.MsgTx, tx stypes.TxOutItem, amount int64, sourceScript []byte, idx int) error {
	sigHashes := btctxscript.NewTxSigHashes(redeemTx)

	var signable btctxscript.Signable
	if tx.VaultPubKey.Equals(c.nodePubKey) {
		signable = btctxscript.NewPrivateKeySignable(c.nodePrivKey)
	} else {
		decodePubKey, err := hex.DecodeString(tx.VaultPubKey.String())
		if err != nil {
			return fmt.Errorf("fail to decode vault public key: %w", err)
		}
		pubKey, err := common.CompressPubKey(decodePubKey)
		if err != nil {
			return fmt.Errorf("fail to compress vault public key: %w", err)
		}
		c.log.Info().Str("pubKey", pubKey).Msg("sign utxo btc")

		signable = newTssSignableBTC(common.PubKey(pubKey), c.tssKeySigner, c.log)
	}

	witness, err := btctxscript.WitnessSignature(redeemTx, sigHashes, idx, amount, sourceScript, btctxscript.SigHashAll, signable, true)
	if err != nil {
		return fmt.Errorf("fail to get witness: %w", err)
	}

	redeemTx.TxIn[idx].Witness = witness
	flag := btctxscript.StandardVerifyFlags
	engine, err := btctxscript.NewEngine(sourceScript, redeemTx, idx, flag, nil, nil, amount)
	if err != nil {
		return fmt.Errorf("fail to create engine: %w", err)
	}
	if err = engine.Execute(); err != nil {
		return fmt.Errorf("fail to execute the script: %w", err)
	}
	return nil
}

func (c *Client) getAddressesFromScriptPubKeyBTC(scriptPubKey btcjson.ScriptPubKeyResult) []string {
	addresses := scriptPubKey.Addresses
	if len(addresses) > 0 {
		return addresses
	}

	if len(scriptPubKey.Hex) == 0 {
		return nil
	}
	buf, err := hex.DecodeString(scriptPubKey.Hex)
	if err != nil {
		c.log.Err(err).Msg("fail to hex decode script pub key")
		return nil
	}
	_, extractedAddresses, _, err := btctxscript.ExtractPkScriptAddrs(buf, c.getChainCfgBTC())
	if err != nil {
		c.log.Err(err).Msg("fail to extract addresses from script pub key")
		return nil
	}
	for _, item := range extractedAddresses {
		addresses = append(addresses, item.String())
	}
	return addresses
}
