package ctx

import (
	"context"
	"io"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/viper"
)

type Context struct {
	FromAddress       string
	ChainID           string
	Codec             codec.Codec
	Input             io.Reader
	Keyring           keyring.Keyring
	KeyringOptions    []keyring.Option
	Output            io.Writer
	OutputFormat      string
	Height            int64
	HomeDir           string
	KeyringDir        string
	From              string
	BroadcastMode     string
	FromName          string
	SignModeStr       string
	UseLedger         bool
	Simulate          bool
	GenerateOnly      bool
	Offline           bool
	SkipConfirm       bool
	NodeURI           string
	Viper             *viper.Viper
	LedgerHasProtobuf bool
	// IsAux is true when the signer is an auxiliary signer (e.g. the tipper).
	IsAux bool
	// TODO: Deprecated (remove).
	LegacyAmino *codec.LegacyAmino
	// CmdContext is the context.Context from the Cobra command.
	CmdContext context.Context
	Client     client.CometRPC
}

func (ctx Context) WithKeyring(k keyring.Keyring) Context {
	ctx.Keyring = k
	return ctx
}

func (ctx Context) WithChainID(chainID string) Context {
	ctx.ChainID = chainID
	return ctx
}

// WithHomeDir returns a copy of the Context with HomeDir set.
func (ctx Context) WithHomeDir(dir string) Context {
	if dir != "" {
		ctx.HomeDir = dir
	}
	return ctx
}

// WithFromName returns a copy of the context with an updated from account name.
func (ctx Context) WithFromName(name string) Context {
	ctx.FromName = name
	return ctx
}

// WithFromAddress returns a copy of the context with an updated from account
// address.
func (ctx Context) WithFromAddress(addr string) Context {
	ctx.FromAddress = addr
	return ctx
}

// WithBroadcastMode returns a copy of the context with an updated broadcast
// mode.
func (ctx Context) WithBroadcastMode(mode string) Context {
	ctx.BroadcastMode = mode
	return ctx
}

// WithNodeURI returns a copy of the context with an updated node URI.
func (ctx Context) WithNodeURI(nodeURI string) Context {
	ctx.NodeURI = nodeURI
	return ctx
}

// WithClient returns a copy of the context with an updated RPC client
// instance.
func (ctx Context) WithClient(client client.CometRPC) Context {
	ctx.Client = client
	return ctx
}
