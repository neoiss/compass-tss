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
	VoteInMethod    = "voteTxIn"
	VoteOutMethod   = "voteTxOut"
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
	maxLogs          int64
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
		maxLogs:          maxLogs,
	}
}

// vaultDepositEvent represent a vault deposit
type vaultDepositEvent struct {
	OrderId ecommon.Hash
	From    ecommon.Address
	Vault   ecommon.Address
	Token   ecommon.Address
	Amount  *big.Int
	To      ecommon.Address
}

func (scp *SmartContractLogParser) parseDeposit(log etypes.Log) (vaultDepositEvent, error) {
	const DepositEventName = "Deposit"
	event := vaultDepositEvent{}
	if err := scp.unpackVaultLog(&event, DepositEventName, log); err != nil {
		return event, fmt.Errorf("fail to unpack deposit event: %w", err)
	}
	return event, nil
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
type swapEvent struct {
	OrderId ecommon.Hash
	From    ecommon.Address // todo not use
	Vault   ecommon.Address
	Token   ecommon.Address
	Amount  *big.Int
	Tochain *big.Int
	To      []byte
	Payload []byte
}

func (scp *SmartContractLogParser) parseSwap(log etypes.Log) (swapEvent, error) {
	const SwapEventName = "Swap"
	event := swapEvent{}
	if err := scp.unpackVaultLog(&event, SwapEventName, log); err != nil {
		return event, fmt.Errorf("fail to unpack swap event: %w", err)
	}
	return event, nil
}

type vaultTransferOutEvent struct {
	OrderId ecommon.Hash
	Vault   ecommon.Address
	Token   ecommon.Address
	Amount  *big.Int
	To      ecommon.Address
	Result  bool
}

func (scp *SmartContractLogParser) parseTransferOut(log etypes.Log) (vaultTransferOutEvent, error) {
	const TransferOutEventName = "TransferOut"
	event := vaultTransferOutEvent{}
	if err := scp.unpackVaultLog(&event, TransferOutEventName, log); err != nil {
		return event, fmt.Errorf("fail to parse transfer out event")
	}
	return event, nil
}

type vaultTransferAllowanceEvent struct {
	OrderId    ecommon.Hash
	FromVault  ecommon.Address
	ToVault    ecommon.Address
	Allowances []*evmTokenAllowance
}

type evmTokenAllowance struct {
	Token  ecommon.Address
	Amount *big.Int
}

func (scp *SmartContractLogParser) parseTransferAllowanceEvent(log etypes.Log) (vaultTransferAllowanceEvent, error) {
	const TransferAllowanceEventName = "TransferAllowance"
	event := vaultTransferAllowanceEvent{}
	if err := scp.unpackVaultLog(&event, TransferAllowanceEventName, log); err != nil {
		return event, fmt.Errorf("fail to parse transfer allowance event")
	}
	return event, nil
}

// THORChainRouterTransferOutAndCall represents a TransferOutAndCall event raised by the THORChainRouter contract.
type THORChainRouterTransferOutAndCall struct {
	Vault        ecommon.Address
	Target       ecommon.Address
	Amount       *big.Int
	FinalAsset   ecommon.Address
	To           ecommon.Address
	AmountOutMin *big.Int
	Memo         string
}

// parseTransferOutAndCall is a log parse operation binding the contract event 0xbda904e26adea40cc083dc36e80fde1641dfdd8b9a035c44022a43e713f73d36.
func (scp *SmartContractLogParser) parseTransferOutAndCall(log etypes.Log) (*THORChainRouterTransferOutAndCall, error) {
	const TransferOutAndCallEventName = "TransferOutAndCall"
	event := new(THORChainRouterTransferOutAndCall)
	if err := scp.unpackVaultLog(event, TransferOutAndCallEventName, log); err != nil {
		return nil, err
	}
	return event, nil
}

func (scp *SmartContractLogParser) GetTxInItem(ll *etypes.Log, txInItem *types.TxInItem) (bool, error) {
	if ll == nil {
		scp.logger.Info().Msg("tx logs are empty return nil")
		return false, nil
	}
	isVaultTransfer := false

	//earlyExit := false
	switch ll.Topics[0].String() {
	case constants.EventOfDeposit.GetTopic().String():
		// router contract , deposit function has re-entrance protection
		depositEvt, err := scp.parseDeposit(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse deposit event")
			return false, err
		}
		txInItem.OrderId = depositEvt.OrderId
		txInItem.To = depositEvt.To.Bytes()
		txInItem.Method = VoteInMethod
		txInItem.TxInType = constants.DEPOSIT
		txInItem.Amount = depositEvt.Amount
		txInItem.Token = depositEvt.Token.Bytes()
		txInItem.Vault = depositEvt.Vault.Bytes()
		txInItem.ToChain = big.NewInt(0) // deposit default zero
		txInItem.From = depositEvt.From.Bytes()

	case constants.EventOfSwap.GetTopic().Hex():
		// it is not legal to have multiple transferOut event , transferOut event should be final
		swapEvt, err := scp.parseSwap(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse swap event")
			return false, err
		}
		txInItem.OrderId = swapEvt.OrderId
		txInItem.From = swapEvt.From.Bytes()
		txInItem.To = swapEvt.To
		txInItem.Method = VoteInMethod
		txInItem.TxInType = constants.SWAP
		txInItem.Amount = swapEvt.Amount
		txInItem.Token = swapEvt.Token.Bytes()
		txInItem.Vault = swapEvt.Vault.Bytes()
		txInItem.ToChain = swapEvt.Tochain
		txInItem.Payload = swapEvt.Payload

	case constants.EventOfTransferOut.GetTopic().Hex():
		// todo convert txOutItem
		// there is no circumstance , router will emit multiple transferAllowance event
		// if that does happen , it means something dodgy happened
		transferOutEvt, err := scp.parseTransferOut(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse transfer out event")
			return false, err
		}
		fmt.Println(" transferOutEvt ------------- ", transferOutEvt)
		//transferOutEvt.To
		//txInItem.OrderId = transferOutEvt.OrderId
		//transferOutEvt.Vault
		//txInItem.Token = transferOutEvt.Token
		//txInItem.Amount = transferOutEvt.Amount
		//transferOutEvt.Result

		isVaultTransfer = false
	case constants.EventOfTransferAllowance.GetTopic().Hex():
		// TODO vault transfer events were only fired by ygg returns
		transferAllowanceEvt, err := scp.parseTransferAllowanceEvent(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse transfer allowance event")
			return false, err
		}
		fmt.Println(" transferAllowanceEvt ------------- ", transferAllowanceEvt)
		//transferAllowanceEvt.Allowances
		//transferAllowanceEvt.FromVault
		//transferAllowanceEvt.ToVault
		//txInItem.OrderId = transferAllowanceEvt.OrderId

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
	case constants.RelayEventOfMigration.GetTopic().String():
		evt, err := scp.parseRelayMigration(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse deposit event")
			return err
		}
		txOutItem.Method = constants.TransferAllowance
		txOutItem.OrderId = evt.OrderId
		txOutItem.Chain = evt.Chain
		txOutItem.FromVault = evt.FromVault
		txOutItem.ToVault = evt.ToVault
		txOutItem.TransactionRate = evt.TransactionRate
		txOutItem.TransactionSize = evt.TransactionSize
		txOutItem.Allowances = make([]types.TokenAllowance, 0)
		for _, ele := range evt.Allowances {
			txOutItem.Allowances = append(txOutItem.Allowances, types.TokenAllowance{
				Token:  ele.Token,
				Amount: ele.Amount,
			})
		}

	case constants.RelayEventOfTransferOut.GetTopic().Hex():
		// it is not legal to have multiple transferOut event , transferOut event should be final
		evt, err := scp.parseRelayTransferOut(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse swap event")
			return err
		}
		txOutItem.Method = constants.TransferOut
		txOutItem.OrderId = evt.OrderId
		txOutItem.Token = evt.Token
		txOutItem.Amount = evt.Amount
		txOutItem.Chain = evt.Chain
		txOutItem.Vault = evt.Vault
		txOutItem.To = evt.To
		txOutItem.TransactionRate = evt.TransactionRate
		txOutItem.TransactionSize = evt.TransactionSize

	case constants.RelayEventOfTransferCall.GetTopic().Hex():
		evt, err := scp.parseRelayTransferCall(*ll)
		if err != nil {
			scp.logger.Err(err).Msg("fail to parse transfer out event")
			return err
		}

		txOutItem.Method = constants.TransferOutCall
		txOutItem.OrderId = evt.OrderId
		txOutItem.Token = evt.Token
		txOutItem.Amount = evt.Amount
		txOutItem.Chain = evt.Chain
		txOutItem.Vault = evt.Vault
		txOutItem.To = evt.To
		txOutItem.Payload = evt.Payload
		txOutItem.TransactionRate = evt.TransactionRate
		txOutItem.TransactionSize = evt.TransactionSize
	}
	return nil
}

type TokenAllowance struct {
	Token  []byte
	Amount *big.Int
}
type RelayMigration struct {
	OrderId         ecommon.Hash
	Chain           *big.Int
	FromVault       []byte
	ToVault         []byte
	Allowances      []TokenAllowance
	TransactionRate *big.Int
	TransactionSize *big.Int
}

func (scp *SmartContractLogParser) parseRelayMigration(log etypes.Log) (RelayMigration, error) {
	const eventName = "Migration"
	event := RelayMigration{}
	if err := scp.unpackVaultLog(&event, eventName, log); err != nil {
		return event, fmt.Errorf("fail to unpack deposit event: %w", err)
	}
	return event, nil
}

type RelayTransferOut struct {
	OrderId         ecommon.Hash
	Token           []byte
	Amount          *big.Int
	Chain           *big.Int
	Vault           []byte
	To              []byte
	TransactionRate *big.Int
	TransactionSize *big.Int
}

func (scp *SmartContractLogParser) parseRelayTransferOut(log etypes.Log) (RelayTransferOut, error) {
	const eventName = "RelayTransferOut"
	event := RelayTransferOut{}
	if err := scp.unpackVaultLog(&event, eventName, log); err != nil {
		return event, fmt.Errorf("fail to unpack deposit event: %w", err)
	}
	return event, nil
}

type RelayTransferCall struct {
	OrderId         ecommon.Hash
	Token           []byte
	Amount          *big.Int
	Chain           *big.Int
	Vault           []byte
	To              []byte
	Payload         []byte
	TransactionRate *big.Int
	TransactionSize *big.Int
}

func (scp *SmartContractLogParser) parseRelayTransferCall(log etypes.Log) (RelayTransferCall, error) {
	const eventName = "RelayTransferCall"
	event := RelayTransferCall{}
	if err := scp.unpackVaultLog(&event, eventName, log); err != nil {
		return event, fmt.Errorf("fail to unpack deposit event: %w", err)
	}
	return event, nil
}
