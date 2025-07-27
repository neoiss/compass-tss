package structure

import (
	"github.com/mapprotocol/compass-tss/constants"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type MaintainerInfo struct {
	Status        uint8          `json:"status"`
	Addr          common.Address `json:"addr"`
	Secp256Pubkey []byte         `json:"secp256Pubkey"`
	Ed25519Pubkey []byte         `json:"ed25519Pubkey"`
	P2pAddress    string         `json:"p2pAddress"`
}

type KeyGen struct {
	Epoch *big.Int         `json:"epoch"`
	Ms    []MaintainerInfo `json:"ms"`
}

type TssPoolParam struct {
	Id        [32]byte
	Epoch     *big.Int
	Pubkey    []byte
	KeyShare  []byte
	Members   []common.Address
	Blames    []common.Address
	Signature []byte
}

type EpochInfo struct {
	Status     uint8
	StartBlock *big.Int
	EndBlock   *big.Int
	Pubkey     []byte
	KeyShare   []byte
	Maitainers []common.Address
}

type VoteTxIn struct {
	TxInType  constants.TxInType
	ToChain   *big.Int
	Height    *big.Int
	FromChain *big.Int
	Amount    *big.Int
	OrderId   [32]byte
	Vault     []byte
	Token     []byte
	From      []byte
	To        []byte
	Payload   []byte
}

type Gas struct {
	Chain          *big.Int
	Pubkey         []byte
	Router         []byte
	GasRate        *big.Int
	TxSize         *big.Int
	TxSizeWithCall *big.Int
}
