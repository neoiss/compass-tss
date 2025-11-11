package evm

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	NativeTokenAddr = "0x0000000000000000000000000000000000000000"
)

type SmartContractLogParser struct {
	logger   zerolog.Logger
	vaultABI *abi.ABI
}

func NewSmartContractLogParser(vaultABI *abi.ABI) SmartContractLogParser {
	return SmartContractLogParser{
		vaultABI: vaultABI,
		logger:   log.Logger.With().Str("module", "SmartContractLogParser").Logger(),
	}
}

// bridgeOutEvent represent a bridge out event
type bridgeOutEvent struct {
	OrderId          ecommon.Hash
	ChainAndGasLimit *big.Int
	TxOutType        uint8
	Vault            []byte
	Token            ecommon.Address
	Amount           *big.Int
	From             ecommon.Address
	RefundAddr       ecommon.Address
	To               []byte
	Payload          []byte
}

func (scp *SmartContractLogParser) parseBridgeOut(log etypes.Log) (*bridgeOutEvent, error) {
	const eventName = "BridgeOut"
	event := bridgeOutEvent{}
	if err := scp.unpackVaultLog(&event, eventName, log); err != nil {
		return nil, fmt.Errorf("fail to unpack BridgeOut event: %w", err)
	}
	return &event, nil
}

func (scp *SmartContractLogParser) unpackVaultLog(out interface{}, event string, log etypes.Log) error {
	if len(log.Topics) == 0 {
		return errors.New("topics field in event log is empty")
	}
	if log.Topics[0] != scp.vaultABI.Events[event].ID {
		return errors.New("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err := scp.vaultABI.UnpackIntoInterface(out, event, log.Data); err != nil {
			return fmt.Errorf("fail to parse event: %w", err)
		}
	}
	var indexed abi.Arguments
	for _, arg := range scp.vaultABI.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return abi.ParseTopics(out, indexed, log.Topics[1:])
}

// swapEvent represent a vault deposit
type bridgeIn struct {
	OrderId          ecommon.Hash
	ChainAndGasLimit *big.Int
	TxInType         uint8
	Vault            []byte
	Sequence         *big.Int
	Sender           ecommon.Address
	Token            ecommon.Address
	Amount           *big.Int
	From             []byte
	To               ecommon.Address
	Data             []byte
}

func (scp *SmartContractLogParser) parseBridgeIn(log etypes.Log) (*bridgeIn, error) {
	const eventName = "BridgeIn"
	event := bridgeIn{}
	if err := scp.unpackVaultLog(&event, eventName, log); err != nil {
		return nil, fmt.Errorf("fail to unpack BridgeIn event: %w", err)
	}
	return &event, nil
}

func (scp *SmartContractLogParser) GetTxInItem(ll *etypes.Log, txInItem *types.TxInItem) error {
	if ll == nil {
		scp.logger.Info().Msg("tx logs are empty return nil")
		return nil
	}

	txInItem.Tx = ll.TxHash.Hex()
	txInItem.LogIndex = ll.Index
	txInItem.Height = big.NewInt(0).SetUint64(ll.BlockNumber)

	switch ll.Topics[0].String() {
	case constants.EventOfBridgeOut.GetTopic().String():
		// router contract , deposit function has re-entrance protection
		evt, err := scp.parseBridgeOut(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse bridge out event")
			return err
		}
		txInItem.Sender = evt.From.Hex()
		txInItem.Amount = evt.Amount
		txInItem.OrderId = evt.OrderId
		txInItem.Token = evt.Token.Bytes()
		txInItem.Vault = evt.Vault
		txInItem.From = evt.From.Bytes()
		txInItem.To = evt.To
		txInItem.Payload = evt.Payload
		txInItem.Method = constants.VoteTxIn
		txInItem.ChainAndGasLimit = evt.ChainAndGasLimit
		txInItem.TxOutType = evt.TxOutType
		txInItem.RefundAddr = evt.RefundAddr.Bytes() // for refund

	case constants.EventOfBridgeIn.GetTopic().Hex():
		// it is not legal to have multiple transferOut event , transferOut event should be final
		evt, err := scp.parseBridgeIn(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse swap event")
			return err
		}
		txInItem.Amount = evt.Amount
		txInItem.OrderId = evt.OrderId
		txInItem.Token = evt.Token.Bytes()
		txInItem.Vault = evt.Vault
		txInItem.From = evt.From
		txInItem.To = evt.To.Bytes()
		txInItem.Method = constants.VoteTxOut
		txInItem.ChainAndGasLimit = evt.ChainAndGasLimit
		txInItem.TxOutType = evt.TxInType
		txInItem.Payload = evt.Data
		txInItem.Sender = evt.Sender.Hex()
		txInItem.Sequence = evt.Sequence
	}

	return nil
}

func (scp *SmartContractLogParser) GetTxOutItem(ll *etypes.Log, txOutItem *types.TxArrayItem) error {
	if ll == nil {
		scp.logger.Info().Msg("Tx logs are empty return nil")
		return nil
	}
	txOutItem.LogIndex = ll.Index // add this
	txOutItem.TxHash = ll.TxHash.String()
	// parse chain and gas limit
	cgl, err := ParseChainAndGasLimit(ll.Topics[2])
	if err != nil {
		return fmt.Errorf("failed to parse chain and gas limit, err: %w", err)
	}
	fmt.Println("GetTxOutItem cgl ", cgl, "txHash", ll.TxHash.String())
	txOutItem.ToChain = cgl.ToChain

	switch ll.Topics[0].String() {
	case constants.EventOfBridgeRelay.GetTopic().Hex(): // oracle
		evt, err := scp.parseBridgeRelay(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse bridgeRelay event")
			return err
		}
		txOutItem.Method = constants.RelaySigned
		txOutItem.OrderId = evt.OrderId
		txOutItem.Token = evt.Token
		txOutItem.Vault = evt.Vault
		txOutItem.To = evt.To
		txOutItem.Amount = evt.Amount
		txOutItem.ChainAndGasLimit = evt.ChainAndGasLimit
		txOutItem.TxType = evt.TxType
		txOutItem.Sequence = evt.Sequence
		txOutItem.From = evt.From
		txOutItem.Data = evt.Data
		txOutItem.Hash = evt.Hash // sign this field
		txOutItem.FromChain = cgl.FromChain

	case constants.EventOfBridgeCompleted.GetTopic().Hex():
		evt, err := scp.parseBridgeCompleted(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse bridgeCompleted event")
			return err
		}
		txOutItem.OrderId = evt.OrderId
		txOutItem.ChainAndGasLimit = evt.ChainAndGasLimit
		txOutItem.Method = constants.BridgeCompleted // todo
		txOutItem.TxType = evt.TxOutType
		txOutItem.Vault = evt.Vault
		txOutItem.Sequence = evt.Sequence
		txOutItem.Sender = evt.Sender.Bytes()
		txOutItem.Data = evt.Data
		txOutItem.FromChain = cgl.FromChain

	case constants.EventOfBridgeRelaySigned.GetTopic().Hex():
		evt, err := scp.parseBridgeRelaySigned(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse bridgeRelaySigned event")
			return err
		}

		txOutItem.Method = constants.BridgeIn
		txOutItem.OrderId = evt.OrderId
		txOutItem.ChainAndGasLimit = evt.ChainAndGasLimit
		txOutItem.Vault = evt.Vault
		txOutItem.Data = evt.RelayData
		txOutItem.Signature = evt.Signature
		txOutItem.FromChain = cgl.FromChain
	default:
		return fmt.Errorf("unknown event topic: %s", ll.Topics[0].String())
	}
	return nil
}

type ChainAndGasLimit struct {
	FromChain *big.Int
	ToChain   *big.Int
	Third     *big.Int
	End       *big.Int
}

func ParseChainAndGasLimit(cgl ecommon.Hash) (*ChainAndGasLimit, error) {
	if cgl.Hex() == constants.ZeroHash {
		return nil, errors.New("chainAndGasLimit is nil")
	}
	bs := cgl.Bytes()
	if len(bs) != 32 {
		return nil, fmt.Errorf("invalid chainAndGasLimit length: %d", len(bs))
	}
	return &ChainAndGasLimit{
		FromChain: big.NewInt(0).SetBytes(bs[0:8]),
		ToChain:   big.NewInt(0).SetBytes(bs[8:16]),
		Third:     big.NewInt(0).SetBytes(bs[16:24]),
		End:       big.NewInt(0).SetBytes(bs[24:32]),
	}, nil
}

type BridgeRelay struct {
	OrderId          ecommon.Hash
	ChainAndGasLimit *big.Int
	TxType           uint8
	Vault            []byte
	To               []byte
	Token            []byte
	Amount           *big.Int
	Sequence         *big.Int
	Hash             [32]byte
	From             []byte
	Data             []byte
}

func (scp *SmartContractLogParser) parseBridgeRelay(log etypes.Log) (*BridgeRelay, error) {
	const eventName = constants.BridgeRelay
	event := BridgeRelay{}
	if err := scp.unpackVaultLog(&event, eventName, log); err != nil {
		return nil, fmt.Errorf("fail to unpack bridge relay event: %w", err)
	}
	return &event, nil
}

type BridgeCompleted struct {
	OrderId          ecommon.Hash
	ChainAndGasLimit *big.Int
	TxOutType        uint8
	Vault            []byte
	Sequence         *big.Int
	Sender           ecommon.Address
	Data             []byte
}

func (scp *SmartContractLogParser) parseBridgeCompleted(log etypes.Log) (*BridgeCompleted, error) {
	const eventName = constants.BridgeCompleted
	event := BridgeCompleted{}
	if err := scp.unpackVaultLog(&event, eventName, log); err != nil {
		return nil, fmt.Errorf("fail to unpack bridge completed event: %w", err)
	}
	return &event, nil
}

type BridgeRelaySigned struct {
	OrderId          ecommon.Hash
	ChainAndGasLimit *big.Int
	Vault            []byte
	RelayData        []byte
	Signature        []byte
}

func (scp *SmartContractLogParser) parseBridgeRelaySigned(log etypes.Log) (*BridgeRelaySigned, error) {
	const eventName = constants.BridgeRelaySigned
	event := BridgeRelaySigned{}
	if err := scp.unpackVaultLog(&event, eventName, log); err != nil {
		return nil, fmt.Errorf("fail to unpack bridge relay signed event: %w", err)
	}
	return &event, nil
}
