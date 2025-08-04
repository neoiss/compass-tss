package mapo

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	ecore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/ethclient"
)

func isAcceptableError(err error) bool {
	return err == nil || err.Error() == txpool.ErrAlreadyKnown.Error() || strings.HasPrefix(err.Error(), ecore.ErrNonceTooLow.Error())
}

// getChainID retrieves the chain id from the node - if this fails we assume local net
func getChainID(client *ethclient.Client, timeout time.Duration) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get chain id, err: %w", err)
	}
	return chainID, err
}
