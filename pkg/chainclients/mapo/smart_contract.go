package mapo

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	//go:embed abi/maintainer.json
	maintainerAbi string
	//go:embed abi/tokenRegister.json
	tokenRegistryABI string
	//go:embed abi/view.json
	viewABI string
	//go:embed abi/relay.json
	relayABI string
	//go:embed abi/tssManager.json
	tssABI string
	//go:embed abi/gasService.json
	gasABI string
)

type bridgeOption func(*Bridge) error

var (
	opts = make([]bridgeOption, 0)
)

func init() {
	opts = append(opts,
		func(b *Bridge) error {
			mainAbi, err := abi.JSON(strings.NewReader(maintainerAbi))
			if err != nil {
				return fmt.Errorf("fail to unmarshal maintainer abi: %w", err)
			}
			b.mainAbi = &mainAbi
			return nil
		},
		func(b *Bridge) error {
			tss, err := abi.JSON(strings.NewReader(tssABI))
			if err != nil {
				return fmt.Errorf("fail to unmarshal tss abi: %w", err)
			}
			b.tssAbi = &tss
			return nil
		},
		func(b *Bridge) error {
			relay, err := abi.JSON(strings.NewReader(relayABI))
			if err != nil {
				return fmt.Errorf("fail to unmarshal relayABI abi: %w", err)
			}
			b.relayAbi = &relay

			return nil
		},
		func(b *Bridge) error {
			gas, err := abi.JSON(strings.NewReader(gasABI))
			if err != nil {
				return fmt.Errorf("fail to unmarshal relayABI abi: %w", err)
			}
			b.gasAbi = &gas

			return nil
		},
		func(b *Bridge) error {
			registry, err := abi.JSON(strings.NewReader(tokenRegistryABI))
			if err != nil {
				return fmt.Errorf("failed to parse token registry abi: %w", err)
			}
			b.tokenRegistry = &registry

			return nil
		},
	)
}

func InitAbi(b *Bridge) error {
	for _, opt := range opts {
		err := opt(b)
		if err != nil {
			return err
		}
	}
	return nil
}
