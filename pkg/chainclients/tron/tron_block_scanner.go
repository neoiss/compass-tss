package tron

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/tron/api"
	"github.com/mr-tron/base58"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var updateGasInterval int64 = 30

type ReportSolvency func(int64) error

type RefBlock struct {
	Timestamp int64
	Height    int64
	Id        string
}

type TronBlockScanner struct {
	cfg                   config.BifrostBlockScannerConfiguration
	logger                zerolog.Logger
	bridge                shareTypes.Bridge
	api                   *api.TronApi
	gatewayAbi            abi.ABI
	refAddress            string // needed for energy estimation via api
	currentFee            uint64
	ethClient             *ethclient.Client
	globalNetworkFeeQueue chan stypes.NetworkFee
}

func NewTronBlockScanner(
	cfg config.BifrostChainConfiguration,
	bridge shareTypes.Bridge,
) (*TronBlockScanner, error) {
	logger := log.Logger.With().
		Str("module", "blockscanner").
		Str("chain", cfg.ChainID.String()).
		Logger()

	ethClient, err := ethclient.Dial(cfg.RPCHost)
	if err != nil {
		return nil, fmt.Errorf("fail to dial ETH rpc host(%s): %w", cfg.RPCHost, err)
	}

	scanner := TronBlockScanner{
		cfg:       cfg.BlockScanner,
		logger:    logger,
		api:       api.NewTronApi(cfg.RPCHost, cfg.BlockScanner.HTTPRequestTimeout),
		bridge:    bridge,
		ethClient: ethClient,
	}
	scanner.gatewayAbi, err = abi.JSON(bytes.NewReader(gatewayABI))
	if err != nil {
		logger.Err(err).Msg("failed to parse ABI")
		return nil, err
	}

	return &scanner, nil
}

func (s *TronBlockScanner) GetHeight() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.HTTPRequestTimeout)
	defer cancel()

	number, err := s.ethClient.BlockNumber(ctx)
	if err != nil {
		s.logger.Err(err).Msg("failed to get latest block")
		return 0, err
	}

	height := int64(number) - ConfirmationBlocks
	if height < 0 {
		height = 0
	}

	return height, nil
}

func (s *TronBlockScanner) FetchMemPool(_ int64) (types.TxIn, error) {
	return types.TxIn{Chain: common.TRONChain}, nil
}

func (s *TronBlockScanner) FetchTxs(
	currentHeight, chainHeight int64,
) (types.TxIn, error) {
	evmContract, _ := s.base58ToHex(s.cfg.Mos) // have 41 prefix
	logs, err := s.ethClient.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(currentHeight),
		ToBlock:   big.NewInt(currentHeight),
		Addresses: []ecommon.Address{ecommon.HexToAddress(strings.Replace(evmContract, "41", "0x", 1))},
		Topics: [][]ecommon.Hash{{
			constants.EventOfBridgeOut.GetTopic(), // txIn -> voteTxIn
			constants.EventOfBridgeIn.GetTopic(),  // txOut -> voteTxOut
		}},
	})
	if err != nil {
		return stypes.TxIn{}, err
	}

	selfId, _ := s.cfg.ChainID.ChainID()
	interval, err := s.bridge.GetMimirWithRef(constants.KeyOfGASFeeGap, selfId.String())
	if err != nil {
		return stypes.TxIn{}, fmt.Errorf("failed to get confirm count: %w", err)
	}
	skip := false
	if interval != 0 && currentHeight%interval != 0 {
		skip = true
	}

	txs, err := s.processTxs(logs)
	if err != nil {
		s.logger.Err(err).Msg("processTxs failed")
		return types.TxIn{}, err
	}

	txIn := types.TxIn{
		Chain:    s.cfg.ChainID,
		TxArray:  txs,
		Filtered: false,
		MemPool:  false,
	}

	if chainHeight-currentHeight > s.cfg.ObservationFlexibilityBlocks {
		return txIn, nil
	}

	if !skip {
		s.updateFees(currentHeight)
	}

	if currentHeight%refBlockInterval != 0 {
		return txIn, nil
	}

	return txIn, nil
}

func (s *TronBlockScanner) GetNetworkFee() (uint64, uint64, uint64) {
	return 1, 1, s.currentFee
}

func (s *TronBlockScanner) base58ToHex(addr string) (string, error) {
	raw, err := base58.Decode(addr)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(raw[:21]), nil
}

// private
// ----------------------------------------------------------------------------
func (s *TronBlockScanner) processTxs(
	logs []etypes.Log,
) ([]*types.TxInItem, error) {
	txInItems := make([]*types.TxInItem, 0)
	cId, _ := s.cfg.ChainID.ChainID()

	for _, ele := range logs {
		// extract the txInItem
		var (
			err error
		)
		// if err := s.blockMetaAccessor.RemoveSignedTxItem(els.TxHash.String()); err != nil {
		// 	s.logger.Err(err).Str("tx hash", els.TxHash.String()).Msg("failed to remove signed tx item")
		// }

		txInItem := &stypes.TxInItem{
			Tx:        ele.TxHash.String()[2:], // drop the "0x" prefix
			LogIndex:  ele.Index,
			Topic:     ele.Topics[0].Hex(),
			Height:    big.NewInt(int64(ele.BlockNumber)),
			FromChain: cId,
		}

		parser := evm.NewSmartContractLogParser(&s.gatewayAbi)
		// txInItem will be changed in parser.getTxInItem function, so if the function return an
		// error txInItem should be abandoned
		tmp := ele
		err = parser.GetTxInItem(&tmp, txInItem)
		if err != nil {
			continue
		}

		// add the txInItem to the txInbound
		txInItems = append(txInItems, txInItem)
	}

	return txInItems, nil
}

func (s *TronBlockScanner) updateFees(height int64) {
	if height%updateGasInterval != 0 {
		return
	}

	params, err := s.api.GetChainParameters()
	if err != nil {
		s.logger.Err(err).Msg("failed get chain parameters")
	}

	var fee, bandwidth, energy int64

	// bandwidth calculation:
	// len(raw_data) + protobuf overhead + max_result_size + signature length
	// len(raw_data) + 3 bytes + 64 bytes + 67 bytes

	// hex data is longer and penalty factor is applied
	// => 576 + 3 + 64 + 67 = 710
	bandwidth = 710 * params.BandwidthFee

	maxEnergy, err := s.getMaxEnergy()
	if err != nil || maxEnergy <= 0 {
		s.logger.Err(err).Msg("failed to get max energy")
		return
	}

	energy = maxEnergy * params.EnergyFee

	// add 1.1 TRX in case the new account needs to be activated:
	// https://developers.tron.network/docs/account#account-activation
	fee = energy + bandwidth + params.MemoFee + 1_100_000

	if fee <= 0 {
		s.logger.Error().Msg("fee is zero")
		return
	}

	// skip sending the network fee if it did not change
	if uint64(fee) == s.currentFee {
		return
	}

	s.currentFee = uint64(fee)

	cId, _ := s.cfg.ChainID.ChainID()
	s.globalNetworkFeeQueue <- stypes.NetworkFee{
		ChainId:             cId,
		Height:              height,
		TransactionSize:     1,
		TransactionSwapSize: 1,
		TransactionRate:     s.currentFee,
	}

	s.logger.Info().
		Int64("height", height).
		Int64("bandwidth", bandwidth).
		Int64("energy", energy).
		Int64("memo_fee", params.MemoFee).
		Int64("energy_fee", params.EnergyFee).
		Int64("bandwidth_fee", params.BandwidthFee).
		Int64("total_fee", fee).
		Msg("updated network fee")
}

func (s *TronBlockScanner) getMaxEnergy() (int64, error) {
	// This value represents the average energy cost of executing a contract transaction
	return 200000, nil
}
