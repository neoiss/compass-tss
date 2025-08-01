package tss

import (
	"time"

	"github.com/mapprotocol/compass-tss/p2p/conversion"
	"github.com/mapprotocol/compass-tss/p2p/messages"
	"github.com/mapprotocol/compass-tss/tss/go-tss/blame"
	"github.com/mapprotocol/compass-tss/tss/go-tss/common"
	"github.com/mapprotocol/compass-tss/tss/go-tss/keygen"
)

func (t *TssServer) Keygen(req keygen.Request) (keygen.Response, error) {
	t.tssKeyGenLocker.Lock()
	defer t.tssKeyGenLocker.Unlock()
	status := common.Success
	msgID, err := t.requestToMsgId(req)
	if err != nil {
		return keygen.Response{}, err
	}

	keygenInstance := keygen.NewTssKeyGen(
		t.p2pCommunication.GetLocalPeerID(),
		t.conf,
		t.localNodePubKey,
		t.p2pCommunication.BroadcastMsgChan,
		t.stopChan,
		t.preParams,
		msgID,
		t.stateManager,
		t.privateKey,
		t.p2pCommunication)

	keygenMsgChannel := keygenInstance.GetTssKeyGenChannels()
	t.p2pCommunication.SetSubscribe(messages.TSSKeyGenMsg, msgID, keygenMsgChannel)
	t.p2pCommunication.SetSubscribe(messages.TSSKeyGenVerMsg, msgID, keygenMsgChannel)
	t.p2pCommunication.SetSubscribe(messages.TSSControlMsg, msgID, keygenMsgChannel)
	t.p2pCommunication.SetSubscribe(messages.TSSTaskDone, msgID, keygenMsgChannel)

	defer func() {
		t.p2pCommunication.CancelSubscribe(messages.TSSKeyGenMsg, msgID)
		t.p2pCommunication.CancelSubscribe(messages.TSSKeyGenVerMsg, msgID)
		t.p2pCommunication.CancelSubscribe(messages.TSSControlMsg, msgID)
		t.p2pCommunication.CancelSubscribe(messages.TSSTaskDone, msgID)

		t.p2pCommunication.ReleaseStream(msgID)
		t.partyCoordinator.ReleaseStream(msgID)
	}()
	sigChan := make(chan string)
	blameMgr := keygenInstance.GetTssCommonStruct().GetBlameMgr()
	joinPartyStartTime := time.Now()
	onlinePeers, leader, errJoinParty := t.joinParty(msgID, req.Version, req.BlockHeight, req.Keys, len(req.Keys)-1, sigChan)
	joinPartyTime := time.Since(joinPartyStartTime)
	if errJoinParty != nil {
		t.logger.Error().Err(errJoinParty).Msgf("failed to joinParty after %s, onlinePeers=%v", joinPartyTime, onlinePeers)

		t.tssMetrics.KeygenJoinParty(joinPartyTime, false)
		t.tssMetrics.UpdateKeyGen(0, false)
		// this indicate we are processing the leaderless join party
		if leader == "NONE" {
			if onlinePeers == nil {
				t.logger.Error().Err(err).Msg("error before we start join party")
				return keygen.Response{
					Status: common.Fail,
					Blame:  blame.NewBlame(blame.InternalError, []blame.Node{}),
				}, nil
			}
			blameNodes, err := blameMgr.NodeSyncBlame(req.Keys, onlinePeers)
			if err != nil {
				t.logger.Err(errJoinParty).Msg("fail to get peers to blame")
			}
			// make sure we blame the leader as well
			t.logger.Error().Err(errJoinParty).Msgf("fail to form keygen party with online:%v", onlinePeers)
			return keygen.Response{
				Status: common.Fail,
				Blame:  blameNodes,
			}, nil

		}

		var blameLeader blame.Blame
		var blameNodes blame.Blame
		blameNodes, err = blameMgr.NodeSyncBlame(req.Keys, onlinePeers)
		if err != nil {
			t.logger.Error().Err(err).Msg("failed to blame nodes for joinParty failure")
		}
		leaderPubKey, err := conversion.GetPubKeyFromPeerIDByEth(leader)
		if err != nil {
			t.logger.Error().Err(err).Msgf("failed to convert peerID->pubkey for leader %s", leader)
			blameLeader = blame.NewBlame(blame.TssSyncFail, []blame.Node{})
		} else {
			blameLeader = blame.NewBlame(blame.TssSyncFail, []blame.Node{{Pubkey: leaderPubKey, BlameData: nil, BlameSignature: nil}})
		}

		if len(onlinePeers) != 0 {
			t.logger.Trace().Msgf("there were %d onlinePeers, adding leader to %d existing nodes blamed",
				len(onlinePeers), len(blameNodes.BlameNodes))
			blameNodes.AddBlameNodes(blameLeader.BlameNodes...)
		} else {
			t.logger.Trace().Msgf("there were %d onlinePeers, setting blame nodes to just the leader",
				len(onlinePeers))
			blameNodes = blameLeader
		}
		t.logger.Error().Err(errJoinParty).Msgf("fail to form keygen party with online:%v", onlinePeers)

		return keygen.Response{
			Status: common.Fail,
			Blame:  blameNodes,
		}, nil

	}

	t.logger.Info().Msg("joinParty succeeded, keygen party formed")
	t.notifyJoinPartyChan()
	t.tssMetrics.KeygenJoinParty(joinPartyTime, true)

	// the statistic of keygen only care about Tss it self, even if the
	// following http response aborts, it still counted as a successful keygen
	// as the Tss model runs successfully.
	beforeKeygen := time.Now()
	k, err := keygenInstance.GenerateNewKey(req)
	keygenTime := time.Since(beforeKeygen)
	if err != nil {
		t.tssMetrics.UpdateKeyGen(keygenTime, false)
		blameNodes := *blameMgr.GetBlame()
		t.logger.Error().Err(err).Msgf("failed to generate key, blaming: %+v", blameNodes.BlameNodes)
		return keygen.NewResponse("", "", common.Fail, blameNodes), err
	} else {
		t.tssMetrics.UpdateKeyGen(keygenTime, true)
	}

	newPubKey, addr, err := conversion.GetTssPubKey(k)
	if err != nil {
		t.logger.Error().Err(err).Msg("failed to generate new tss pubkey from generated key")
		status = common.Fail
	}

	blameNodes := *blameMgr.GetBlame()
	t.logger.Trace().Msgf("returning from keygen with status=%d, blaming=%+v", status, blameNodes.BlameNodes)
	return keygen.NewResponse(
		newPubKey,
		addr.String(),
		status,
		blameNodes,
	), nil
}
