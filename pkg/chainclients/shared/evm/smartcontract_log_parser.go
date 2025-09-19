package evm

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	NativeTokenAddr = "0x0000000000000000000000000000000000000000"
)

type (
	contractAddressValidator func(addr *ecommon.Address, includeWhitelist bool) bool
	assetResolver            func(token string) (common.Asset, error)
	tokenDecimalResolver     func(token string) int64
	amountConverter          func(token string, amt *big.Int) cosmos.Uint
)

type SmartContractLogParser struct {
	addressValidator contractAddressValidator
	assetResolver    assetResolver
	decimalResolver  tokenDecimalResolver
	amtConverter     amountConverter
	logger           zerolog.Logger
	vaultABI         *abi.ABI
	nativeAsset      common.Asset
}

func NewSmartContractLogParser(validator contractAddressValidator,
	resolver assetResolver,
	decimalResolver tokenDecimalResolver,
	amtConverter amountConverter,
	vaultABI *abi.ABI,
	nativeAsset common.Asset,
	maxLogs int64,
) SmartContractLogParser {
	return SmartContractLogParser{
		addressValidator: validator,
		assetResolver:    resolver,
		decimalResolver:  decimalResolver,
		vaultABI:         vaultABI,
		amtConverter:     amtConverter,
		logger:           log.Logger.With().Str("module", "SmartContractLogParser").Logger(),
		nativeAsset:      nativeAsset,
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
func (scp *SmartContractLogParser) GetTxInItem(ll *etypes.Log, txInItem *types.TxInItem) (bool, error) {
	if ll == nil {
		scp.logger.Info().Msg("tx logs are empty return nil")
		return false, nil
	}
	isVaultTransfer := false
	txInItem.Tx = ll.TxHash.Hex()
	txInItem.LogIndex = ll.Index
	txInItem.Height = big.NewInt(0).SetUint64(ll.BlockNumber)

	//earlyExit := false
	switch ll.Topics[0].String() {
	case constants.EventOfBridgeOut.GetTopic().String():
		// router contract , deposit function has re-entrance protection
		evt, err := scp.parseBridgeOut(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse bridge out event")
			return false, err
		}
		// txInItem.FromChain = evt.FromChain
		// txInItem.ToChain = evt.ToChain
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
		// txInItem.Memo
		// txInItem.ObservedVaultPubKey = evt.ObservedVaultPubKey
		// txInItem.GasUsed = evt.GasUsed

	case constants.EventOfBridgeIn.GetTopic().Hex():
		// it is not legal to have multiple transferOut event , transferOut event should be final
		evt, err := scp.parseBridgeIn(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse swap event")
			return false, err
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

	return isVaultTransfer, nil
}

func (scp *SmartContractLogParser) GetTxOutItem(ll *etypes.Log, txOutItem *types.TxArrayItem) error {
	if ll == nil {
		scp.logger.Info().Msg("Tx logs are empty return nil")
		return nil
	}
	txOutItem.LogIndex = ll.Index // add this
	txOutItem.TxHash = ll.TxHash.String()

	switch ll.Topics[0].String() {
	case constants.EventOfBridgeRelay.GetTopic().Hex():
		evt, err := scp.parseBridgeRelay(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse deposit event")
			return err
		}
		txOutItem.Method = constants.TransferAllowance
		txOutItem.OrderId = evt.OrderId
		txOutItem.Token = evt.Token
		txOutItem.Vault = evt.Vault
		txOutItem.To = evt.To
		txOutItem.Amount = evt.Amount
		// txOutItem.Chain = evt.Chain // todo handler
		txOutItem.ChainAndGasLimit = evt.ChainAndGasLimit
		txOutItem.TxOutType = evt.TxOutType
		txOutItem.Sequence = evt.Sequence
		txOutItem.From = evt.From
		txOutItem.Data = evt.Data

	case constants.EventOfBridgeCompleted.GetTopic().Hex():
		// todo handler
	}
	return nil
}

type BridgeRelay struct {
	OrderId          ecommon.Hash
	ChainAndGasLimit *big.Int
	TxOutType        uint8
	Vault            []byte
	Sequence         *big.Int
	Token            []byte
	Amount           *big.Int
	From             []byte
	To               []byte
	Data             []byte
}

func (scp *SmartContractLogParser) parseBridgeRelay(log etypes.Log) (*BridgeRelay, error) {
	const eventName = constants.EventOfBridgeRelay
	event := BridgeRelay{}
	if err := scp.unpackVaultLog(&event, eventName.String(), log); err != nil {
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
	Sender           common.Address
	Data             []byte
}

func (scp *SmartContractLogParser) parseBridgeCompleted(log etypes.Log) (*BridgeCompleted, error) {
	const eventName = constants.EventOfBridgeCompleted
	event := BridgeCompleted{}
	if err := scp.unpackVaultLog(&event, eventName.String(), log); err != nil {
		return nil, fmt.Errorf("fail to unpack bridge completed event: %w", err)
	}
	return &event, nil
}
