package utxo

import (
	"fmt"

	dogeec "github.com/eager7/dogd/btcec"
	dogechaincfg "github.com/eager7/dogd/chaincfg"
	dogewire "github.com/eager7/dogd/wire"
	dogetxscript "github.com/mapprotocol/compass-tss/txscript/dogd-txscript"

	"github.com/mapprotocol/compass-tss/common"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
)

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
