package structure

import (
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
