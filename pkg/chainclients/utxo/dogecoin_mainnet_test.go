//go:build !testnet
// +build !testnet

package utxo

import (
	"math/big"

	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/mapclient/types"
)

func (s *DogecoinSuite) TestGetAddress(c *C) {
	pubkey := common.PubKey("thorpub1addwnpepqt7qug8vk9r3saw8n4r803ydj2g3dqwx0mvq5akhnze86fc536xcy2cr8a2")
	addr := s.client.GetAddress(pubkey)
	c.Assert(addr, Equals, "DCdSuatdjCqdWJFB6LEeFweabLiypVxLsz")
}

func (s *DogecoinSuite) TestConfirmationCountReady(c *C) {
	c.Assert(s.client.ConfirmationCountReady(types.TxIn{
		Chain:    common.DOGEChain,
		TxArray:  nil,
		Filtered: true,
		MemPool:  false,
	}), Equals, true)

	c.Assert(s.client.ConfirmationCountReady(types.TxIn{
		Chain: common.DOGEChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(2),
				Tx:     "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo: "MEMO",
			},
		},
		Filtered: true,
		MemPool:  true,
	}), Equals, true)
	s.client.currentBlockHeight.Store(3)
	c.Assert(s.client.ConfirmationCountReady(types.TxIn{
		Chain: common.DOGEChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(2),
				Tx:     "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo: "MEMO",
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 0,
	}), Equals, true)

	c.Assert(s.client.ConfirmationCountReady(types.TxIn{
		Chain: common.DOGEChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(2),
				Tx:     "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo: "MEMO",
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 5,
	}), Equals, false)
}

func (s *DogecoinSuite) TestGetConfirmationCount(c *C) {

	// no tx in item , confirmation count should be 0
	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain:   common.DOGEChain,
		TxArray: nil,
	}), Equals, int64(0))
	// mempool txin , confirmation count should be 0
	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain: common.DOGEChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(2),
				Tx:     "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo: "MEMO",
			},
		},
		Filtered:             true,
		MemPool:              true,
		ConfirmationRequired: 0,
	}), Equals, int64(0))

	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain: common.DOGEChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(2),
				Tx:     "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo: "MEMO",
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 0,
	}), Equals, int64(0))

	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain: common.DOGEChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(2),
				Tx:     "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo: "MEMO",
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 0,
	}), Equals, int64(0))

	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain: common.DOGEChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(2),
				Tx:     "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo: "MEMO",
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 0,
	}), Equals, int64(0))

	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain: common.DOGEChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(2),
				Tx:     "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo: "MEMO",
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 0,
	}), Equals, int64(20))
}
