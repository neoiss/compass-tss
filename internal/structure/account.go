package structure

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type MaintainerInfo struct {
	Status            uint8          `json:"status,omitempty"`
	Account           common.Address `json:"account,omitempty"`
	LastHeartbeatTime *big.Int       `json:"last_heartbeat_time,omitempty"`
	LastActiveEpoch   *big.Int       `json:"last_active_epoch,omitempty"`
	Secp256Pubkey     []byte         `json:"secp_256_pubkey,omitempty"`
	Ed25519Pubkey     []byte         `json:"ed_25519_pubkey,omitempty"`
	P2pAddress        string         `json:"p_2_p_address,omitempty"`
}

type KeyGen struct {
	Epoch *big.Int         `json:"epoch"`
	Ms    []MaintainerInfo `json:"ms"`
}

type TssPoolParam struct {
	Epoch     *big.Int
	Pubkey    []byte
	KeyShare  []byte
	Members   []common.Address
	Blames    []common.Address
	Signature []byte
}

type EpochInfo struct {
	ElectedBlock  uint64
	StartBlock    uint64
	EndBlock      uint64
	MigratedBlock uint64
	Maintainers   []common.Address
}

type TxInItem struct {
	TxInType         uint8
	OrderId          [32]byte
	ChainAndGasLimit *big.Int
	Height           *big.Int
	Token            []byte
	Amount           *big.Int
	From             []byte
	Vault            []byte
	To               []byte
	Payload          []byte
}

type VoteTxOut struct {
	TxOutType        uint8
	OrderId          [32]byte
	ChainAndGasLimit *big.Int
	Height           *big.Int
	GasUsed          *big.Int
	Sequence         *big.Int
	Amount           *big.Int
	Sender           common.Address
	Token            []byte
	From             []byte
	To               []byte
	Vault            []byte
	Data             []byte
}

type Gas struct {
	Chain          *big.Int
	Pubkey         []byte
	Router         []byte
	GasRate        *big.Int
	TxSize         *big.Int
	TxSizeWithCall *big.Int
}
