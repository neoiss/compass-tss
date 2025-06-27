package utxo

import (
	"fmt"

	bchec "github.com/gcash/bchd/bchec"
	bchchaincfg "github.com/gcash/bchd/chaincfg"
	bchwire "github.com/gcash/bchd/wire"
	bchtxscript "github.com/mapprotocol/compass-tss/txscript/bchd-txscript"

	"github.com/mapprotocol/compass-tss/common"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
)

func (c *Client) getChainCfgBCH() *bchchaincfg.Params {
	switch common.CurrentChainNetwork {
	case common.TestNet:
		return &bchchaincfg.TestNet3Params
	case common.MainNet:
		return &bchchaincfg.MainNetParams
	default:
		c.log.Fatal().Msg("unsupported network")
		return nil
	}
}

func (c *Client) signUTXOBCH(redeemTx *bchwire.MsgTx, tx stypes.TxOutItem, amount int64, sourceScript []byte, idx int) error {
	var signable bchtxscript.Signable
	if tx.VaultPubKey.Equals(c.nodePubKey) {
		signable = bchtxscript.NewPrivateKeySignable((*bchec.PrivateKey)(c.nodePrivKey))
	} else {
		signable = newTssSignableBCH(tx.VaultPubKey, c.tssKeySigner, c.log)
	}

	sig, err := bchtxscript.RawTxInECDSASignature(redeemTx, idx, sourceScript, bchtxscript.SigHashAll, signable, amount)
	if err != nil {
		return fmt.Errorf("fail to get witness: %w", err)
	}

	pkData := signable.GetPubKey().SerializeCompressed()
	sigscript, err := bchtxscript.NewScriptBuilder().AddData(sig).AddData(pkData).Script()
	if err != nil {
		return fmt.Errorf("fail to build signature script: %w", err)
	}
	redeemTx.TxIn[idx].SignatureScript = sigscript
	flag := bchtxscript.StandardVerifyFlags
	engine, err := bchtxscript.NewEngine(sourceScript, redeemTx, idx, flag, nil, nil, amount)
	if err != nil {
		return fmt.Errorf("fail to create engine: %w", err)
	}
	if err = engine.Execute(); err != nil {
		return fmt.Errorf("fail to execute the script: %w", err)
	}
	return nil
}
