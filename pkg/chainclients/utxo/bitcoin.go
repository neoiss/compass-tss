package utxo

import (
	"encoding/hex"
	"fmt"

	btcjson "github.com/btcsuite/btcd/btcjson"
	btcchaincfg "github.com/btcsuite/btcd/chaincfg"
	btcwire "github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	btctxscript "github.com/mapprotocol/compass-tss/txscript/txscript"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/pkg/address"
)

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
		c.log.Info().Str("pubKey", tx.VaultPubKey.String()).Msg("sign utxo btc")
		signable = newTssSignableBTC(tx.VaultPubKey, c.tssKeySigner, c.log)
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

func (c *Client) decodeBTCAddress(toAddress string, isMigrate bool, txHash string) (btcutil.Address, bool, error) {
	if isMigrate {
		addr, err := btcutil.DecodeAddress(toAddress, c.getChainCfgBTC())
		if err != nil {
			c.log.Error().Err(err).Str("relayHash", txHash).Str("toAddress", toAddress).Msg("fail to decode bitcoin address")
			return nil, false, fmt.Errorf("fail to decode bitcoin address: %w", err)
		}
		return addr, false, nil
	}

	addr, err := address.DecodeBitcoinAddress(toAddress, c.getChainCfgBTC())
	if err == nil {
		return addr, false, nil
	}

	c.log.Error().Err(err).Str("relayHash", txHash).Str("toAddress", toAddress).Msg("fail to decode bitcoin address")

	chainID, _ := c.cfg.ChainID.ChainID()
	defaultAddress, err := c.bridge.GetMimirWithBytes(constants.KeyOfTransferFailedReceiver, chainID.String())
	if err != nil {
		c.log.Error().Err(err).Str("relayHash", txHash).Str("chain", c.cfg.ChainID.String()).Msg("fail to get default receiver")
		return nil, false, fmt.Errorf("fail to get default address config: %w", err)
	}
	addr, err = address.DecodeBitcoinAddress(hex.EncodeToString(defaultAddress), c.getChainCfgBTC())
	if err != nil {
		c.log.Error().Err(err).Str("relayHash", txHash).Str("addr", hex.EncodeToString(defaultAddress)).Msg("fail to decode bitcoin address")
		return nil, false, fmt.Errorf("fail to decode next address: %w", err)
	}
	return addr, true, nil
}
