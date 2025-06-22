package tss

import (
	"fmt"
	"github.com/mapprotocol/compass-tss/internal/keys"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/mapprotocol/compass-tss/tss/go-tss/keygen"
	"github.com/mapprotocol/compass-tss/tss/go-tss/tss"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
	"github.com/mapprotocol/compass-tss/x/types"
)

// KeyGen is
type KeyGen struct {
	keys           *keys.Keys
	logger         zerolog.Logger
	client         *http.Client
	server         *tss.TssServer
	bridge         shareTypes.Bridge
	currentVersion string
	lastCheck      time.Time
}

// NewTssKeyGen create a new instance of TssKeyGen which will look after TSS key stuff
func NewTssKeyGen(keys *keys.Keys, server *tss.TssServer, bridge shareTypes.Bridge) (*KeyGen, error) {
	if keys == nil {
		return nil, fmt.Errorf("keys is nil")
	}
	return &KeyGen{
		keys:   keys,
		logger: log.With().Str("module", "tss_keygen").Logger(),
		client: &http.Client{
			Timeout: time.Second * 130,
		},
		server: server,
		bridge: bridge,
	}, nil
}

func (kg *KeyGen) getVersion() string {
	requestTime := time.Now()
	//if !kg.currentVersion.Equals(semver.Version{}) && requestTime.Sub(kg.lastCheck).Seconds() < constants.MAPRelayChainBlockTime.Seconds() {
	//	return kg.currentVersion
	//}
	version, err := kg.bridge.GetThorchainVersion()
	if err != nil {
		kg.logger.Err(err).Msg("fail to get current thorchain version")
		return kg.currentVersion
	}
	kg.currentVersion = version
	kg.lastCheck = requestTime
	return kg.currentVersion
}

func (kg *KeyGen) GenerateNewKey(keygenBlockHeight int64, pKeys common.PubKeys) (pk common.PubKeySet, blame types.Blame, err error) {
	// No need to do key gen
	if len(pKeys) == 0 {
		return common.EmptyPubKeySet, types.Blame{}, nil
	}

	// add some logging
	defer func() {
		if blame.IsEmpty() {
			kg.logger.Info().Int64("height", keygenBlockHeight).Str("pubkey", pk.String()).Msg("tss keygen results success")
		} else {
			blames := make([]string, len(blame.BlameNodes))
			for i := range blame.BlameNodes {
				var pk common.PubKey
				pk, err = common.NewPubKey(blame.BlameNodes[i].Pubkey)
				if err != nil {
					kg.logger.Error().Err(err).Int64("height", keygenBlockHeight).Str("pubkey", blame.BlameNodes[i].Pubkey).Msg("tss keygen results error")
					continue
				}
				var acc cosmos.AccAddress
				acc, err = pk.GetThorAddress()
				if err != nil {
					kg.logger.Error().Err(err).Int64("height", keygenBlockHeight).Str("pubkey", pk.String()).Msg("tss keygen results error")
					continue
				}
				blames[i] = acc.String()
			}
			sort.Strings(blames)
			kg.logger.Info().Int64("height", keygenBlockHeight).Str("pubkey", pk.String()).Str("round", blame.Round).Str("blames", strings.Join(blames, ", ")).Str("reason", blame.FailReason).Msg("tss keygen results blame")
		}
	}()

	var keys []string
	for _, item := range pKeys {
		keys = append(keys, item.String())
	}
	keyGenReq := keygen.Request{
		Keys: keys,
	}
	//currentVersion := kg.getVersion()
	keyGenReq.Version = kg.getVersion()

	// Use the churn try's block to choose the same leader for every node in an Asgard,
	// since a successful keygen requires every node in the Asgard to take part.
	keyGenReq.BlockHeight = keygenBlockHeight

	ch := make(chan bool, 1)
	defer close(ch)
	timer := time.NewTimer(30 * time.Minute)
	defer timer.Stop()

	fmt.Println("processKeygenBlock GenerateNewKey 222222 -------------- ")
	var resp keygen.Response
	go func() {
		resp, err = kg.server.Keygen(keyGenReq)
		ch <- true
	}()

	select {
	case <-ch:
		// do nothing
	case <-timer.C:
		panic("tss keygen timeout")
	}

	// copy blame to our own struct
	blame = types.Blame{
		FailReason: resp.Blame.FailReason,
		IsUnicast:  resp.Blame.IsUnicast,
		Round:      resp.Blame.Round,
		BlameNodes: make([]types.Node, len(resp.Blame.BlameNodes)),
	}
	for i, n := range resp.Blame.BlameNodes {
		blame.BlameNodes[i].Pubkey = n.Pubkey
		blame.BlameNodes[i].BlameData = n.BlameData
		blame.BlameNodes[i].BlameSignature = n.BlameSignature
	}
	fmt.Println("processKeygenBlock GenerateNewKey 333333 -------------- ")
	if err != nil {
		// the resp from kg.server.Keygen will not be nil
		if blame.IsEmpty() {
			blame.FailReason = err.Error()
		}
		return common.EmptyPubKeySet, blame, fmt.Errorf("fail to keygen,err:%w", err)
	}

	fmt.Println("processKeygenBlock GenerateNewKey 444444 -------------- ", resp.PubKey)
	cpk, err := common.NewPubKey(resp.PubKey)
	if err != nil {
		return common.EmptyPubKeySet, blame, fmt.Errorf("fail to create common.PubKey,%w", err)
	}

	// TODO later on THORNode need to have both secp256k1 key and ed25519
	return common.NewPubKeySet(cpk, cpk), blame, nil
}
