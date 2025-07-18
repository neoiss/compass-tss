package utxo

import (
	"fmt"

	ltcec "github.com/ltcsuite/ltcd/btcec"
	ltcchaincfg "github.com/ltcsuite/ltcd/chaincfg"
	ltcwire "github.com/ltcsuite/ltcd/wire"
	ltctxscript "github.com/mapprotocol/compass-tss/txscript/ltcd-txscript"

	"github.com/mapprotocol/compass-tss/common"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
)

func (c *Client) getChainCfgLTC() *ltcchaincfg.Params {
	cn := common.CurrentChainNetwork
	switch cn {
	case common.TestNet:
		return &ltcchaincfg.TestNet4Params
	case common.MainNet:
		return &ltcchaincfg.MainNetParams
	}
	return nil
}

func (c *Client) signUTXOLTC(redeemTx *ltcwire.MsgTx, tx stypes.TxOutItem, amount int64, sourceScript []byte, idx int) error {
	sigHashes := ltctxscript.NewTxSigHashes(redeemTx)

	var signable ltctxscript.Signable
	if tx.VaultPubKey.Equals(c.nodePubKey) {
		signable = ltctxscript.NewPrivateKeySignable((*ltcec.PrivateKey)(c.nodePrivKey))
	} else {
		signable = newTssSignableLTC(tx.VaultPubKey, c.tssKeySigner, c.log)
	}

	witness, err := ltctxscript.WitnessSignature(redeemTx, sigHashes, idx, amount, sourceScript, ltctxscript.SigHashAll, signable, true)
	if err != nil {
		return fmt.Errorf("fail to get witness: %w", err)
	}

	redeemTx.TxIn[idx].Witness = witness
	flag := ltctxscript.StandardVerifyFlags
	engine, err := ltctxscript.NewEngine(sourceScript, redeemTx, idx, flag, nil, nil, amount)
	if err != nil {
		return fmt.Errorf("fail to create engine: %w", err)
	}
	if err = engine.Execute(); err != nil {
		return fmt.Errorf("fail to execute the script: %w", err)
	}
	return nil
}
