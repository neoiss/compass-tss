package mapo

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/internal/structure"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	stypes "github.com/mapprotocol/compass-tss/x/types"
	"github.com/pkg/errors"
)

func (b *Bridge) Register() error {
	method := constants.Register
	priv, err := b.keys.GetPrivateKey()
	if err != nil {
		return err
	}

	ethPrivateKey, err := evm.GetPrivateKey(priv)
	if err != nil {
		return err
	}
	pk, err := keys.PublicKeyFromPrivate(ethPrivateKey)
	if err != nil {
		return err
	}
	publicKeyBytes := ecrypto.FromECDSAPub(pk)
	// fmt.Println("publicKeyBytes ----------- ", publicKeyBytes)
	// fmt.Println("publicKeyBytes ----------- ", ecommon.Bytes2Hex(publicKeyBytes))

	input, err := b.mainAbi.Pack(method, publicKeyBytes[1:], publicKeyBytes[1:], b.cfg.Addr)
	if err != nil {
		return errors.Wrap(err, "fail to pack input")
	}

	ctx := context.TODO()
	txBytes, err := b.assemblyTx(ctx, input, 0, b.cfg.Maintainer)
	if err != nil {
		return errors.Wrap(err, "fail to assembly tx")
	}
	txId, err := b.Broadcast(txBytes)
	if err != nil {
		return errors.Wrap(err, "fail to broadcast tx")
	}
	b.logger.Info().Str("txId", txId).Msg("Register maintainer successfully")

	return nil
}

// GetKeygenBlock retrieves keygen request for the given block height from mapBridge
func (b *Bridge) GetKeygenBlock() (*structure.KeyGen, error) {
	method := constants.ElectionEpoch
	input, err := b.mainAbi.Pack(method)
	if err != nil {
		return nil, errors.Wrap(err, "fail to pack input")
	}
	var epoch *big.Int
	err = b.callContract(&epoch, b.cfg.Maintainer, method, input, b.mainAbi)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call contract")
	}

	b.logger.Info().Any("epoch", epoch).Msg("GetKeygenBlock-----------")
	if epoch.Uint64() == 0 { // not in epoch
		return nil, nil
	}
	tssStatus, err := b.getTssStatus(epoch)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get tss status")
	}
	b.logger.Info().Any("epoch", epoch).Any("status", tssStatus).Msg("The epoch info")
	switch tssStatus {
	case constants.TssStatusPending, constants.TssStatusConsensus:
		break
	case constants.TssStatusUnknown:
		return nil, nil
	case constants.TssStatusFailed:
		return nil, fmt.Errorf("tss status (%d)failed", tssStatus)
	default:
		b.epoch = epoch
		// b.epoch = big.NewInt(0)
		b.logger.Info().Any("epoch", epoch).Any("tssStatus", tssStatus).
			Msg("The epoch tss status is completed")
		return nil, nil
	}

	b.logger.Info().Int64("epoch", epoch.Int64()).Msg("KeyGen Block")
	// done
	ret, err := b.GetEpochInfo(epoch)
	if err != nil {
		return nil, err
	}
	ms, err := b.GetNodeAccounts(ret.Maintainers)
	if err != nil {
		return nil, err
	}

	idx := -1
	selfAddr, _ := b.keys.GetEthAddress()
	members := make([]ecommon.Address, 0)
	for i, ele := range ms {
		members = append(members, ele.Account)
		if strings.EqualFold(ele.Account.Hex(), selfAddr.Hex()) {
			idx = i
		}
	}
	if idx == -1 {
		b.epoch = epoch
		b.logger.Debug().Any("self", selfAddr).Any("elect", ms).Msg("This node is not in the election period")
		return nil, nil
	}

	currEpochHash, err := b.genHash(epoch, members)
	if err != nil {
		return nil, err
	}
	// chech hash
	if strings.EqualFold(b.epochHash.Hex(), currEpochHash.Hex()) {
		b.logger.Info().Any("epoch", epoch).Any("history", b.epochHash.Hex()).
			Any("curEpoch", currEpochHash).Msg("The epoch is same")
		return nil, nil
	}

	// use compressed pk, dont modify this code
	for idx, item := range ms {
		pks := make([]byte, 0)
		pks = append(pks, 4)
		pks = append(pks, item.Secp256Pubkey...)
		epk, err := ecrypto.UnmarshalPubkey(pks)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal ECDSA public key: %w", err)
		}
		pubBytes := ecrypto.CompressPubkey(epk)
		ms[idx].Secp256Pubkey = pubBytes
	}
	return &structure.KeyGen{
		Epoch: epoch,
		Ms:    ms,
	}, nil
}

func (b *Bridge) GetEpochInfo(epoch *big.Int) (*structure.EpochInfo, error) {
	method := constants.GetEpochInfo
	input, err := b.mainAbi.Pack(method, epoch)
	if err != nil {
		return nil, err
	}

	ret := struct {
		Info structure.EpochInfo
	}{}
	err = b.callContract(&ret, b.cfg.Maintainer, method, input, b.mainAbi)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to call %s", method)
	}

	return &ret.Info, nil
}

// GetNodeAccount retrieves node account for this address from mapBridge
func (b *Bridge) GetNodeAccount(addr string) (*structure.MaintainerInfo, error) {
	method := constants.GetMaintainerInfos
	input, err := b.mainAbi.Pack(method,
		[]ecommon.Address{ecommon.HexToAddress(addr)})
	if err != nil {
		return nil, err
	}

	type Back struct {
		Infos []structure.MaintainerInfo `json:"infos"`
	}
	var ret Back
	err = b.callContract(&ret, b.cfg.Maintainer, method, input, b.mainAbi)
	if err != nil {
		return nil, err
	}
	return &ret.Infos[0], nil
}

// GetNodeAccounts retrieves all node accounts from mapBridge
func (b *Bridge) GetNodeAccounts(addrs []ecommon.Address) ([]structure.MaintainerInfo, error) {
	method := constants.GetMaintainerInfos
	input, err := b.mainAbi.Pack(method, addrs)
	if err != nil {
		return nil, err
	}

	type Back struct {
		Infos []structure.MaintainerInfo `json:"infos"`
	}
	var ret Back
	err = b.callContract(&ret, b.cfg.Maintainer, method, input, b.mainAbi)
	if err != nil {
		return nil, err
	}

	return ret.Infos, nil
}

func (b *Bridge) callContract(ret interface{}, addr, method string, input []byte, abi *abi.ABI) error {
	to := ecommon.HexToAddress(addr)
	outPut, err := b.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		From: constants.ZeroAddress,
		To:   &to,
		Data: input,
	}, nil)
	if err != nil {
		return errors.Wrapf(err, "unable to call contract %s", method)
	}

	outputs := abi.Methods[method].Outputs
	unpack, err := outputs.Unpack(outPut)
	if err != nil {
		return errors.Wrap(err, "unpack output")
	}

	if err = outputs.Copy(ret, unpack); err != nil {
		return errors.Wrap(err, "copy output")
	}
	return nil
}

// FetchNodeStatus get current node status from mapBridge
func (b *Bridge) FetchNodeStatus() (stypes.NodeStatus, error) {
	addr, err := b.keys.GetEthAddress()
	if err != nil {
		return stypes.NodeStatus_Unknown, nil
	}
	// done
	na, err := b.GetNodeAccount(addr.String())
	if err != nil {
		return stypes.NodeStatus_Unknown, fmt.Errorf("failed to get node status: %w", err)
	}

	return stypes.NodeStatus(na.Status), nil
}

func (b *Bridge) FetchActiveNodes() ([]common.PubKey, error) {
	// done
	ret, err := b.GetEpochInfo(b.epoch)
	if err != nil {
		return nil, err
	}
	na, err := b.GetNodeAccounts(ret.Maintainers)
	if err != nil {
		return nil, fmt.Errorf("fail to get node accounts: %w", err)
	}
	active := make([]common.PubKey, 0)
	for _, item := range na {
		if stypes.NodeStatus(item.Status) == stypes.NodeStatus_Active {
			active = append(active, common.PubKey(ecommon.Bytes2Hex(item.Secp256Pubkey)))
		}
	}
	return active, nil
}

func (b *Bridge) genHash(epoch *big.Int, members []ecommon.Address) (ecommon.Hash, error) {
	data := []interface{}{epoch, len(members), members}

	encoded, err := rlp.EncodeToBytes(data)
	if err != nil {
		return ecommon.Hash{}, fmt.Errorf("RLP failed: %v", err)
	}
	hash := ecommon.BytesToHash(ecrypto.Keccak256(encoded))
	return hash, nil

}
