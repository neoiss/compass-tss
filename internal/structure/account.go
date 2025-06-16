package structure

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type MaintainerInfo struct {
	Status        uint8          `json:"status"`
	Version       *big.Int       `json:"version"`
	Addr          common.Address `json:"addr"`
	Secp256Pubkey []byte         `json:"secp256Pubkey"`
	Ed25519Pubkey []byte         `json:"ed25519Pubkey"`
	P2pAddress    string         `json:"p2pAddress"`
}

type KeyGen struct {
	Epoch *big.Int          `json:"epoch"`
	Ms    []*MaintainerInfo `json:"ms"`
}

type TssPoolParam struct {
	Id        [32]byte         //id = keccak256(abi.encodePacked(pubkey, members, epoch, blames, chainIds));
	Epoch     *big.Int         // epoch 第一次就先填1。 合约是有选举的 ，选举后是会弹出一个日志。日志有下一轮的 的epoch
	Pubkey    []byte           // 生成的 tss地址的公钥
	Members   []common.Address // 由哪些节点一起生成的这个tss
	Blames    []common.Address // 如果再参与生成 tss地址的时候 有被选上的但是 通信没有反应的列入到惩罚里面
	Signature []byte           // tss地址的公钥对应的私钥 对上面的id的签名
}
