package utxo

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/mapprotocol/compass-tss/constants"

	dogeec "github.com/eager7/dogd/btcec"
	"github.com/eager7/dogd/chaincfg"
	dogechaincfg "github.com/eager7/dogd/chaincfg"
	dogewire "github.com/eager7/dogd/wire"
	"github.com/eager7/dogutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	dogetxscript "github.com/mapprotocol/compass-tss/txscript/dogd-txscript"

	"github.com/mapprotocol/compass-tss/common"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
)

func EncodeDOGEAddress(addr dogutil.Address) (string, error) {
	script := hex.EncodeToString(addr.ScriptAddress())

	switch addr.(type) {
	case *dogutil.AddressWitnessPubKeyHash: // P2WPKH
		return P2WPKH + script, nil
	case *dogutil.AddressWitnessScriptHash: // P2WSH
		return P2WSH + script, nil
	case *dogutil.AddressPubKeyHash: // P2PKH
		return P2PKH + script, nil
	case *dogutil.AddressScriptHash: // P2SH
		return P2SH + script, nil
	default:
		return "", errUnknownAddressType
	}
}

func DecodeDOGEAddress(addr string, network *chaincfg.Params) (dogutil.Address, error) {
	if !common.HasHexPrefix(addr) {
		addr = "0x" + addr
	}
	if len(addr) < 4 {
		return nil, errors.New("invalid address")
	}
	prefix := addr[:4]
	publicKey := addr[4:]
	publicKeyLen := len(publicKey)
	publicKeyBytes := ethcommon.Hex2Bytes(publicKey)

	switch prefix {
	case P2WPKHOrP2WSH:
		switch publicKeyLen {
		case 40: // P2WPKH
			return dogutil.NewAddressWitnessPubKeyHash(publicKeyBytes, network)
		case 64: // P2WSH
			return dogutil.NewAddressWitnessScriptHash(publicKeyBytes, network)
		default:
			return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
		}
	case P2PKHOrP2TR:
		switch publicKeyLen {
		case 40: // P2PKH
			return dogutil.NewAddressPubKeyHash(publicKeyBytes, network) // base 58
		//case 64: // P2TR dogecoin not support taproot
		default:
			return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
		}

	case P2SH:
		return dogutil.NewAddressScriptHashFromHash(publicKeyBytes, network) // base 58
	default:
		return nil, newUnsupportedPublicKeyLenError(prefix, publicKeyLen)
	}
}

func (c *Client) getChainCfgDOGE() *dogechaincfg.Params {
	switch common.CurrentChainNetwork {
	case common.TestNet:
		return &dogechaincfg.TestNet3Params
	case common.MainNet:
		return &dogechaincfg.MainNetParams
	default:
		c.log.Fatal().Msg("unsupported network")
		return nil
	}
}

func (c *Client) signUTXODOGE(redeemTx *dogewire.MsgTx, tx stypes.TxOutItem, amount int64, sourceScript []byte, idx int) error {
	var signable dogetxscript.Signable
	if tx.VaultPubKey.Equals(c.nodePubKey) {
		signable = dogetxscript.NewPrivateKeySignable((*dogeec.PrivateKey)(c.nodePrivKey))
	} else {
		c.log.Info().Str("pubKey", tx.VaultPubKey.String()).Msg("sign utxo doge")
		signable = newTssSignableDOGE(tx.VaultPubKey, c.tssKeySigner, c.log)
	}

	sig, err := dogetxscript.RawTxInSignature(redeemTx, idx, sourceScript, dogetxscript.SigHashAll, signable)
	if err != nil {
		return fmt.Errorf("fail to get witness: %w", err)
	}

	pkData := signable.GetPubKey().SerializeCompressed()
	sigscript, err := dogetxscript.NewScriptBuilder().AddData(sig).AddData(pkData).Script()
	if err != nil {
		return fmt.Errorf("fail to build signature script: %w", err)
	}
	redeemTx.TxIn[idx].SignatureScript = sigscript
	flag := dogetxscript.StandardVerifyFlags
	engine, err := dogetxscript.NewEngine(sourceScript, redeemTx, idx, flag, nil, nil, amount)
	if err != nil {
		return fmt.Errorf("fail to create engine: %w", err)
	}
	if err = engine.Execute(); err != nil {
		return fmt.Errorf("fail to execute the script: %w", err)
	}
	return nil
}

func (c *Client) decodeDOGEAddress(toAddress string, isMigrate bool, txHash string) (dogutil.Address, bool, error) {
	if isMigrate {
		addr, err := dogutil.DecodeAddress(toAddress, c.getChainCfgDOGE())
		if err != nil {
			c.log.Error().Err(err).Str("relayHash", txHash).Str("toAddress", toAddress).Msg("fail to decode dogecoin address")
			return nil, false, fmt.Errorf("fail to decode dogecoin address: %w", err)
		}
		return addr, false, nil
	}

	addr, err := DecodeDOGEAddress(toAddress, c.getChainCfgDOGE())
	if err == nil {
		return addr, false, nil
	}

	c.log.Error().Err(err).Str("relayHash", txHash).Str("toAddress", toAddress).Msg("fail to decode dogecoin address")
	defaultAddress, err := c.bridge.GetMimirWithBytes(constants.KeyOfTransferFailedReceiver, c.cfg.ChainID.String())
	if err != nil {
		c.log.Error().Err(err).Str("relayHash", txHash).Str("chain", c.cfg.ChainID.String()).Msg("fail to get default receiver")
		return nil, false, fmt.Errorf("fail to get default address config: %w", err)
	}
	addr, err = DecodeDOGEAddress(hex.EncodeToString(defaultAddress), c.getChainCfgDOGE())
	if err != nil {
		c.log.Error().Err(err).Str("relayHash", txHash).Str("addr", hex.EncodeToString(defaultAddress)).Msg("fail to decode dogecoin address")
		return nil, false, fmt.Errorf("fail to decode next address: %w", err)
	}
	return addr, true, nil
}
