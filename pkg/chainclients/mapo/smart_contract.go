package mapo

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	idPackAbi = `[{"inputs":[{"internalType":"bytes","name":"pubkey","type":"bytes"},{"internalType":"address[]","name":"members","type":"address[]"},{"internalType":"uint256","name":"epoch","type":"uint256"},{"internalType":"address[]","name":"blames","type":"address[]"},{"internalType":"bytes","name":"keyShares","type":"bytes"}],"name":"idPack","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

var (
	//go:embed abi/maintainer.json
	maintainerAbi string
	//go:embed abi/tokenRegister.json
	tokenRegistryABI string
)

func newMaintainerABi() (*abi.ABI, error) {
	maintainer, err := abi.JSON(strings.NewReader(maintainerAbi))
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal maintainer abi: %w", err)
	}

	return &maintainer, nil
}

func NewIdABi() (*abi.ABI, error) {
	maintainer, err := abi.JSON(strings.NewReader(idPackAbi))
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal maintainer abi: %w", err)
	}

	return &maintainer, nil
}

func NewTokenRegistry() (*abi.ABI, error) {
	registry, err := abi.JSON(strings.NewReader(tokenRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token registry abi: %w", err)
	}

	return &registry, nil
}
