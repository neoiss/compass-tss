package utxo

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/eager7/dogutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	btypes "github.com/mapprotocol/compass-tss/blockscanner/types"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/utxo"
	mem "github.com/mapprotocol/compass-tss/x/memo"
)

type TxProcessor struct {
	cli *Client
}

func NewTxProcessor(client *Client) *TxProcessor {
	return &TxProcessor{cli: client}
}

type TxProcessContext struct {
	Tx                 *btcjson.TxRawResult
	Height             int64
	VinZeroTxs         map[string]*btcjson.TxRawResult
	Sender             string
	Memo               string
	ParsedMemo         mem.Memo
	ToBytes            []byte
	Refund             bool
	ChainID            *big.Int
	Output             btcjson.Vout
	Receiver           string
	NativeTokenAddress []byte
}

func (p *TxProcessor) isValidTx(tx *btcjson.TxRawResult, height int64, isMemPool bool) bool {
	if p.cli.ignoreTx(tx, height) {
		p.cli.log.Debug().Int64("height", height).Str("txid", tx.Txid).Msg("ignore tx not matching format")
		return false
	}
	// RBF enabled transaction will not be observed until committed to block
	if p.cli.isRBFEnabled(tx) && isMemPool {
		p.cli.log.Debug().Int64("height", height).Str("txid", tx.Txid).Msg("ignore RBF enabled tx in mempool")
		return false
	}
	return true
}

func (p *TxProcessor) PrepareContext(tx *btcjson.TxRawResult, height int64, isMemPool bool, vinZeroTxs map[string]*btcjson.TxRawResult) (*TxProcessContext, error) {
	if !p.isValidTx(tx, height, isMemPool) {
		return nil, nil
	}

	sender, err := p.cli.getSender(tx, vinZeroTxs)
	if err != nil {
		return nil, fmt.Errorf("fail to get sender from tx: %w", err)
	}

	memo, err := p.cli.getMemo(tx)
	if err != nil {
		return nil, fmt.Errorf("fail to get memo from tx: %w", err)
	}

	if len([]byte(memo)) > constants.MaxMemoSize {
		return nil, fmt.Errorf("memo (%s) longer than max allow length (%d)", memo, constants.MaxMemoSize)
	}

	parsedMemo, toBytes, refund, err := p.ParseMemo(memo, tx.Txid)
	if err != nil {
		return nil, err
	}

	chainID, err := p.cli.GetChain().ChainID()
	if err != nil {
		return nil, fmt.Errorf("fail to get chain id: %w, chain: %s", err, p.cli.GetChain())
	}

	output, err := p.cli.getOutput(sender, tx, false)
	if err != nil {
		if errors.Is(err, btypes.ErrFailOutputMatchCriteria) {
			p.cli.log.Debug().Int64("height", height).Str("txid", tx.Txid).Msg("ignore tx not matching format")
			return nil, nil
		}
		return nil, fmt.Errorf("fail to get output from tx: %w", err)
	}

	addresses := p.cli.getAddressesFromScriptPubKey(output.ScriptPubKey)
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no addresses found in output")
	}
	receiver := addresses[0]

	if p.cli.cfg.ChainID.Equals(common.BCHChain) {
		receiver = p.cli.stripBCHAddress(receiver)
	}

	nativeToken, ok := p.cli.GetChain().NativeToken()
	if !ok {
		return nil, fmt.Errorf("fail to get native token, chain: %s", p.cli.GetChain())
	}
	nativeTokenAddress, err := p.cli.bridge.GetTokenAddress(chainID, nativeToken)
	if err != nil {
		return nil, fmt.Errorf("fail to get token address: %w, chainID: %s, token: %s", err, chainID, nativeToken)
	}

	return &TxProcessContext{
		Tx:                 tx,
		Height:             height,
		VinZeroTxs:         vinZeroTxs,
		Sender:             sender,
		Memo:               memo,
		ParsedMemo:         parsedMemo,
		ToBytes:            toBytes,
		Refund:             refund,
		ChainID:            chainID,
		Output:             output,
		Receiver:           receiver,
		NativeTokenAddress: nativeTokenAddress,
	}, nil
}

func (p *TxProcessor) ParseMemo(memo string, txid string) (mem.Memo, []byte, bool, error) {
	var (
		refund  bool
		toBytes []byte
	)
	parsedMemo, err := mem.ParseMemo(memo)
	if err != nil {
		p.cli.log.Debug().Err(err).Str("txid", txid).Str("memo", memo).Msg("fail to parse memo")
		refund = true
	} else {
		toBytes, err = parsedMemo.GetChain().DecodeAddress(parsedMemo.GetDestination())
		if err != nil {
			p.cli.log.Debug().Err(err).Str("txid", txid).Str("memo", memo).Msg("fail to decode memo")
			refund = true
		}
	}
	return parsedMemo, toBytes, refund, nil
}

func (p *TxProcessor) ProcessAsgardSenderTx(ctx *TxProcessContext) (types.TxInItem, error) {
	vaultPbuKey, err := utxo.GetAsgardPubKeyByAddress(p.cli.cfg.ChainID, p.cli.bridge, common.Address(ctx.Sender))
	if err != nil {
		return types.TxInItem{}, fmt.Errorf("fail to get vault pub key by address: %w", err)
	}

	amount, err := btcutil.NewAmount(ctx.Tx.Vout[0].Value)
	if err != nil {
		return types.TxInItem{}, fmt.Errorf("fail to parse amount(%f): %w", ctx.Tx.Vout[0].Value, err)
	}
	txResult, err := p.cli.rpc.GetRawTransactionVerboseWithFee(ctx.Tx.Txid)
	if err != nil {
		return types.TxInItem{}, fmt.Errorf("fail to get tx result: %w", err)
	}
	fee, err := btcutil.NewAmount(txResult.Fee)
	if err != nil {
		return types.TxInItem{}, fmt.Errorf("fail to parse amount(%f): %w", ctx.Tx.Vout[0].Value, err)
	}

	toBytes, txOutType, payload, err := p.ProcessMemoTypeForAsgardSender(ctx.Receiver, ctx.ParsedMemo)
	if err != nil {
		return types.TxInItem{}, err
	}

	chainAndGasLimit := make([]byte, 32)
	toChain := ethcommon.LeftPadBytes(ctx.ChainID.Bytes(), 8)
	copy(chainAndGasLimit[8:16], toChain)

	txIn := types.TxInItem{
		Tx:               ctx.Tx.Txid,
		Memo:             ctx.Memo,
		Height:           new(big.Int).SetInt64(ctx.Height),
		Amount:           big.NewInt(int64(amount)),
		OrderId:          ethcommon.HexToHash(ctx.ParsedMemo.GetOrderID()),
		GasUsed:          big.NewInt(int64(fee)),
		Token:            ctx.NativeTokenAddress,
		Vault:            vaultPbuKey,
		From:             nil,
		To:               toBytes,
		Payload:          payload,
		Method:           constants.VoteTxOut,
		LogIndex:         0,
		ChainAndGasLimit: new(big.Int).SetBytes(chainAndGasLimit),
		TxOutType:        uint8(txOutType),
		Sequence:         big.NewInt(0),
		Topic:            constants.EventOfBridgeIn.GetTopic().String(),
		Timestamp:        ctx.Tx.Blocktime,
	}
	p.cli.log.Info().Int64("height", ctx.Height).Str("txid", ctx.Tx.Txid).Interface("txIn", txIn).Msg("got tx in")
	return txIn, nil
}

func (p *TxProcessor) ProcessMemoTypeForAsgardSender(toAddr string, parsedMemo mem.Memo) ([]byte, constants.TxInType, []byte, error) {
	var (
		toBytes   []byte
		txOutType constants.TxInType
		payload   []byte
		err       error
	)

	switch parsedMemo.GetType() {
	case mem.TxInbound:
		toBytes, err = p.EncodeAddressToBytes(toAddr)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("fail to encode address: %w", err)
		}
		txOutType = constants.TRANSFER
		payload, err = utxo.EncodePayload(nil, nil, nil)
	case mem.TxMigrate:
		// when the tx type is migration, to bytes must be empty.
		toBytes = []byte{}
		txOutType = constants.MIGRATE
		payload, err = utxo.GetAsgardPubKeyByAddress(p.cli.cfg.ChainID, p.cli.bridge, common.Address(toAddr))
		if err != nil {
			return nil, 0, nil, fmt.Errorf("fail to get vault pub key by address: %w", err)
		}
	case mem.TxRefund:
		toBytes, err = p.EncodeAddressToBytes(toAddr)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("fail to encode address: %w", err)
		}
		txOutType = constants.REFUND
		payload, err = utxo.EncodePayload(nil, nil, nil)
	default:
		return nil, 0, nil, fmt.Errorf("unsupported tx type: %s", parsedMemo.GetType())
	}

	if err != nil {
		return nil, 0, nil, err
	}

	return toBytes, txOutType, payload, nil
}

func (p *TxProcessor) ProcessAsgardReceiverTx(ctx *TxProcessContext) (types.TxInItem, error) {
	if !p.cli.isValidUTXO(ctx.Output.ScriptPubKey.Hex) {
		return types.TxInItem{}, fmt.Errorf("invalid utxo")
	}

	amount, err := btcutil.NewAmount(ctx.Output.Value)
	if err != nil {
		return types.TxInItem{}, fmt.Errorf("fail to parse float64: %w, value: %f", err, ctx.Output.Value)
	}
	amt := uint64(amount.ToUnit(btcutil.AmountSatoshi))

	nativeToken, ok := p.cli.GetChain().NativeToken()
	if !ok {
		return types.TxInItem{}, fmt.Errorf("fail to get native token, chain: %s", p.cli.GetChain())
	}

	mapChainID, err := common.MAPChain.ChainID()
	if err != nil {
		return types.TxInItem{}, fmt.Errorf("fail to get chain id: %w, chain: %s", err, common.MAPChain)
	}

	txOutType, destChainID, payload, err := p.ProcessMemoTypeForAsgardReceiver(ctx.ParsedMemo, ctx.Refund, nativeToken, mapChainID, ctx.ToBytes)
	if err != nil {
		return types.TxInItem{}, err
	}

	fromBytes, err := p.EncodeAddressToBytes(ctx.Sender)
	if err != nil {
		return types.TxInItem{}, fmt.Errorf("fail to encode sender address: %w", err)
	}

	pubKey, err := utxo.GetAsgardPubKeyByAddress(p.cli.cfg.ChainID, p.cli.bridge, common.Address(ctx.Receiver))
	if err != nil {
		return types.TxInItem{}, fmt.Errorf("fail to get vault address2pubkey mapped: %w", err)
	}

	chainAndGasLimit := make([]byte, 32)
	fromChain := ethcommon.LeftPadBytes(ctx.ChainID.Bytes(), 8)
	toChain := ethcommon.LeftPadBytes(destChainID.Bytes(), 8)

	toBytes := ctx.ToBytes
	if destChainID.Cmp(mapChainID) != 0 {
		destChainID = mapChainID
		toChain = ethcommon.LeftPadBytes(mapChainID.Bytes(), 8)
		toBytes = p.cli.bridge.GetFusionReceiver().Bytes()
	}
	copy(chainAndGasLimit[0:8], fromChain)
	copy(chainAndGasLimit[8:16], toChain)

	txIn := types.TxInItem{
		Tx:               ctx.Tx.Txid,
		Memo:             ctx.Memo,
		Sender:           ctx.Sender,
		FromChain:        ctx.ChainID,
		ToChain:          destChainID,
		Height:           big.NewInt(ctx.Height),
		Amount:           new(big.Int).SetUint64(amt),
		OrderId:          generateOrderID(ctx.ChainID.String(), ctx.Tx.Txid),
		GasUsed:          big.NewInt(0),
		Token:            ctx.NativeTokenAddress,
		Vault:            pubKey,
		From:             fromBytes,
		To:               toBytes,
		Payload:          payload,
		Method:           constants.VoteTxIn,
		LogIndex:         0,
		ChainAndGasLimit: new(big.Int).SetBytes(chainAndGasLimit),
		TxOutType:        uint8(txOutType),
		RefundAddr:       fromBytes,
		Topic:            constants.EventOfBridgeOut.GetTopic().String(),
		Timestamp:        ctx.Tx.Blocktime,
	}
	p.cli.log.Info().Int64("height", ctx.Height).Str("txid", ctx.Tx.Txid).Interface("txIn", txIn).Msg("got tx in")
	return txIn, nil
}

func (p *TxProcessor) ProcessMemoTypeForAsgardReceiver(parsedMemo mem.Memo, refund bool, nativeToken string, mapChainID *big.Int, toBytes []byte) (constants.TxInType, *big.Int, []byte, error) {
	var (
		txOutType   constants.TxInType
		destChainID = big.NewInt(0)
		payload     []byte
		err         error
	)

	if refund {
		txOutType = constants.TRANSFER
		destChainID = mapChainID
		payload, err = utxo.EncodePayload(nil, nil, nil)
		if err != nil {
			return 0, nil, nil, fmt.Errorf("fail to encode empty payload: %w", err)
		}
		return txOutType, destChainID, payload, nil
	}

	switch parsedMemo.GetType() {
	case mem.TxOutbound:
		destToken := parsedMemo.GetToken()
		txOutType = constants.TRANSFER
		destChainID, err = p.cli.bridge.GetChainID(parsedMemo.GetChain().String())
		if err != nil {
			return 0, nil, nil, fmt.Errorf("fail to get destination chain id: %w, chain: %s", err, parsedMemo.GetChain())
		}
		payload, err = p.cli.encodePayload(nativeToken, destToken, mapChainID, destChainID, toBytes, parsedMemo)
		if err != nil {
			// todo refund
			return 0, nil, nil, fmt.Errorf("fail to encode payload: %w", err)
		}
	case mem.TxAdd:
		txOutType = constants.DEPOSIT
		destChainID = mapChainID
		payload, err = utxo.EncodePayload(nil, nil, nil)
		if err != nil {
			return 0, nil, nil, fmt.Errorf("fail to encode empty payload: %w", err)
		}
	default:
		return 0, nil, nil, fmt.Errorf("unsupported tx type: %s", parsedMemo.GetType())
	}

	return txOutType, destChainID, payload, nil
}

func (p *TxProcessor) EncodeAddressToBytes(addr string) ([]byte, error) {
	var (
		addrStr string
		err     error
	)

	switch p.cli.cfg.ChainID {
	case common.BTCChain:
		var address btcutil.Address
		address, err = btcutil.DecodeAddress(addr, p.cli.getChainCfgBTC())
		if err != nil {
			return nil, fmt.Errorf("fail to decode bitcoin address(%s): %w", addr, err)
		}
		addrStr, err = EncodeBitcoinAddress(address)
		if err != nil {
			return nil, fmt.Errorf("fail to encode bitcoin address(%s): %w", address.String(), err)
		}
	case common.DOGEChain:
		var address dogutil.Address
		address, err = dogutil.DecodeAddress(addr, p.cli.getChainCfgDOGE())
		if err != nil {
			return nil, fmt.Errorf("fail to decode doge address(%s): %w", addr, err)
		}
		addrStr, err = EncodeDOGEAddress(address)
		if err != nil {
			return nil, fmt.Errorf("fail to encode doge address(%s): %w", address.String(), err)
		}
	default:
		return nil, fmt.Errorf("unsupported chain: %s", p.cli.cfg.ChainID)
	}

	toBytes, err := hex.DecodeString(strings.TrimPrefix(addrStr, "0x"))
	if err != nil {
		return nil, fmt.Errorf("fail to decode hex address(%s): %w", addrStr, err)
	}

	return toBytes, nil
}
