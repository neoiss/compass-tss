package tss

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/itchio/lzma"
	"github.com/rs/zerolog/log"

	"github.com/mapprotocol/compass-tss/constants"
	sharedTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/x/types"
)

func RecoverKeyShares(mapchain sharedTypes.Bridge) error {
	tctx := mapchain.GetContext()

	// fetch the node account
	na, err := mapchain.GetNodeAccount(tctx.FromAddress)
	if err != nil {
		return fmt.Errorf("fail to get node account: %w", err)
	}

	// skip recovery if the current node is not active
	if types.NodeStatus(na.Status) != types.NodeStatus_Active {
		log.Info().Msgf("%s is not active, skipping key shares recovery", na.Addr)
		return nil
	}

	// get the current epoch info
	epochInfo, err := mapchain.GetEpochInfo(big.NewInt(0))
	if err != nil {
		return fmt.Errorf("fail to get current epoch info: %w", err)
	}
	// todo
	vault := common.Bytes2Hex(epochInfo.Pubkey)
	keysharesPath := filepath.Join(constants.DefaultHome, fmt.Sprintf("localstate-%s.json", vault))

	keyShare, err := mapchain.GetKeyShare()
	if err != nil {
		return fmt.Errorf("fail to get current epoch info: %w", err)
	}

	// skip recovery if key shares for the nodes current vault already exist
	if _, err = os.Stat(keysharesPath); !os.IsNotExist(err) {
		log.Info().Msgf("key shares for %s already exist, skipping recovery", vault)
		return nil
	}

	if err := recoverKeyShares(keysharesPath, keyShare, os.Getenv("SIGNER_SEED_PHRASE")); err != nil {
		return err
	}
	// success
	log.Info().Str("path", keysharesPath).Msgf("recovered key shares for %s", na.Addr)
	return nil
}

func recoverKeyShares(path string, keyShares []byte, passphrase string) error {
	// open key shares file
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open keyshares file: %w", err)
	}
	defer f.Close()

	// decrypt and decompress into place
	var decrypted []byte
	decrypted, err = DecryptKeyShares(keyShares, passphrase)
	if err != nil {
		return fmt.Errorf("failed to decrypt key shares: %w", err)
	}
	cmpDec := lzma.NewReader(bytes.NewReader(decrypted))
	if _, err = io.Copy(f, cmpDec); err != nil {
		return fmt.Errorf("failed to decompress key shares: %w", err)
	}

	return nil
}
