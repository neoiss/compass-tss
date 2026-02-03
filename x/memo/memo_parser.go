package memo

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

const maxNameLength = 10      // todo
var maxBps = big.NewInt(5000) // todo

type parser struct {
	memo           string
	txType         TxType
	parts          []string
	errs           []error
	requiredFields int
}

// func newParser(ctx cosmos.Context, keeper keeper.Keeper, version semver.Version, memo string) (parser, error) {
func newParser(memo string) (parser, error) {
	if len(memo) == 0 {
		return parser{}, fmt.Errorf("memo can't be empty")
	}

	parts := strings.Split(memo, "|")
	memoType, err := StringToTxType(parts[0])
	if err != nil {
		return parser{}, err
	}

	return parser{
		memo:   memo,
		txType: memoType,
		parts:  parts,
		errs:   make([]error, 0),
	}, nil
}

func (p *parser) getType() TxType {
	return p.txType
}

func (p *parser) incRequired(required bool) {
	if required {
		p.requiredFields += 1
	}
}

func (p *parser) parse() (mem Memo, err error) {
	defer func() {
		if err == nil {
			err = p.Error()
		}
	}()
	switch p.getType() {
	case TxAdd:
		return p.ParseAddLiquidityMemo()
	case TxInbound:
		return p.ParseInboundMemo()
	case TxOutbound:
		return p.ParseOutboundMemo()
	case TxRefund:
		return p.ParseRefundMemo()
	case TxMigrate:
		return p.ParseMigrateMemo()
	default:
		return EmptyMemo, fmt.Errorf("TxType not supported: %s", p.getType().String())
	}
}

func (p *parser) addErr(err error) {
	p.errs = append(p.errs, err)
}

func (p *parser) Error() error {
	p.hasMinParams(p.requiredFields + 1)
	if len(p.errs) == 0 {
		return nil
	}
	errStrs := make([]string, len(p.errs))
	for i, err := range p.errs {
		errStrs[i] = err.Error()
	}
	err := fmt.Errorf("MEMO: %s\nPARSE FAILURE(S): %s", p.memo, strings.Join(errStrs, "-"))
	return err
}

// check if memo has enough parameters
func (p *parser) hasMinParams(count int) {
	if len(p.parts) < count {
		p.addErr(fmt.Errorf("not enough parameters: %d/%d", len(p.parts), count))
	}
}

// Safe accessor for split memo parts - always returns empty
// string for indices that are out of bounds.
func (p *parser) get(idx int) string {
	if idx < 0 || len(p.parts) <= idx {
		return ""
	}
	return p.parts[idx]
}

func (p *parser) getMinAmount(idx int) *big.Int {
	return p.parseMinAmount(p.get(idx))
}

func (p *parser) getAffiliates(idx int) []*Affiliate {
	return p.parseAffiliates(p.get(idx))
}

// Safe accessor for split memo parts - always returns empty string for indices that are
// out of bounds. Returns the sub-index of a split part (with separator "/").
func (p *parser) getSubIndex(idx, subIdx int) string {
	if idx < 0 || len(p.parts) <= idx {
		return ""
	}
	subParts := strings.Split(p.parts[idx], "/")
	if subIdx < 0 || len(subParts) <= subIdx {
		return ""
	}
	return subParts[subIdx]
}

func (p *parser) getInt64(idx int, required bool, def int64) int64 {
	p.incRequired(required)
	value, err := strconv.ParseInt(p.get(idx), 10, 64)
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an int64: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

func (p *parser) getUint(idx int, required bool, def uint64) cosmos.Uint {
	p.incRequired(required)
	value, err := cosmos.ParseUint(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an uint: %w", p.get(idx), err))
		}
		return cosmos.NewUint(def)
	}
	return value
}

func (p *parser) getUintWithScientificNotation(idx int, required bool, def uint64) cosmos.Uint {
	p.incRequired(required)
	f, _, err := big.ParseFloat(p.get(idx), 10, 0, big.ToZero)
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an uint with sci notation: %w", p.get(idx), err))
		}
		return cosmos.NewUint(def)
	}
	i := new(big.Int)
	f.Int(i) // Note: fractional part will be discarded
	result := cosmos.NewUintFromBigInt(i)
	return result
}

func (p *parser) getUintWithMaxValue(idx int, required bool, def, max uint64) cosmos.Uint {
	value := p.getUint(idx, required, def)
	if value.GT(cosmos.NewUint(max)) {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("%s cannot exceed '%d'", p.get(idx), max))
		}
		return cosmos.NewUint(max)
	}
	return value
}

func (p *parser) getAccAddress(idx int, required bool, def cosmos.AccAddress) cosmos.AccAddress {
	p.incRequired(required)
	value, err := cosmos.AccAddressFromBech32(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an AccAddress: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

func (p *parser) getAddress(idx int, required bool, def common.Address) common.Address {
	p.incRequired(required)
	value, err := common.NewAddress(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an Address: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

//func (p *parser) getAddressWithKeeper(idx int, required bool, def common.Address, chain common.Chain) common.Address {
//	p.incRequired(required)
//	if p.keeper == nil {
//		return p.getAddress(2, required, common.NoAddress)
//	}
//	addr, err := FetchAddress(p.ctx, p.keeper, p.get(idx), chain)
//	if err != nil {
//		if required || p.get(idx) != "" {
//			p.addErr(fmt.Errorf("cannot parse '%s' as an Address: %w", p.get(idx), err))
//		}
//	}
//	return addr
//}

func (p *parser) getStringArrayBySeparator(idx int, required bool, separator string) []string {
	p.incRequired(required)
	value := p.get(idx)
	if value == "" {
		return []string{}
	}
	return strings.Split(value, separator)
}

func (p *parser) getUintArrayBySeparator(idx int, required bool, separator string) []cosmos.Uint {
	p.incRequired(required)
	value := p.get(idx)
	if value == "" {
		return []cosmos.Uint{}
	}
	strArray := strings.Split(value, separator)
	result := make([]cosmos.Uint, 0, len(strArray))
	for _, str := range strArray {
		u, err := cosmos.ParseUint(str)
		if err != nil {
			if required || str != "" {
				p.addErr(fmt.Errorf("cannot parse '%s' as an uint: %w", str, err))
			}
			return []cosmos.Uint{}
		}
		result = append(result, u)
	}
	return result
}

//func (p *parser) getAddressAndRefundAddressWithKeeper(idx int, required bool, def common.Address, chain common.Chain) (common.Address, common.Address) {
//	p.incRequired(required)
//
//	//nolint:ineffassign
//	destination := common.NoAddress
//	refundAddress := common.NoAddress
//	addresses := p.get(idx)
//
//	if strings.Contains(addresses, "/") {
//		parts := strings.SplitN(addresses, "/", 2)
//		if p.keeper == nil {
//			dest, err := common.NewAddress(parts[0])
//			if err != nil {
//				if required || parts[0] != "" {
//					p.addErr(fmt.Errorf("cannot parse '%s' as an Address: %w", parts[0], err))
//				}
//			}
//			destination = dest
//		} else {
//			destination = p.getAddressFromString(parts[0], chain, required)
//		}
//		if len(parts) > 1 {
//			refundAddress, _ = common.NewAddress(parts[1])
//		}
//	} else {
//		destination = p.getAddressWithKeeper(idx, false, common.NoAddress, chain)
//	}
//
//	if destination.IsEmpty() && !refundAddress.IsEmpty() {
//		p.addErr(fmt.Errorf("refund address is set but destination address is empty"))
//	}
//
//	return destination, refundAddress
//}

//func (p *parser) getAddressFromString(val string, chain common.Chain, required bool) common.Address {
//	addr, err := FetchAddress(p.ctx, p.keeper, val, chain)
//	if err != nil {
//		if required || val != "" {
//			p.addErr(fmt.Errorf("cannot parse '%s' as an Address: %w", val, err))
//		}
//	}
//	return addr
//}

func (p *parser) getChain(idx int, required bool, def common.Chain) common.Chain {
	p.incRequired(required)
	value, err := common.NewChain(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as a chain: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

//func (p *parser) getAsset(idx int, required bool, def common.Asset) common.Asset {
//	p.incRequired(required)
//	value, err := common.NewAssetWithShortCodes(p.version, p.get(idx))
//	if err != nil && (required || p.get(idx) != "") {
//		p.addErr(fmt.Errorf("cannot parse '%s' as an asset: %w", p.get(idx), err))
//		return def
//	}
//	return value
//}

func (p *parser) getTxID(idx int, required bool, def common.TxID) common.TxID {
	p.incRequired(required)
	value, err := common.NewTxID(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as tx hash: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

func (p *parser) getBase64Bytes(idx int, required bool, def []byte) []byte {
	p.incRequired(required)
	value, err := base64.StdEncoding.DecodeString(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as a base64 string: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

func (p *parser) parseMinAmount(s string) *big.Int {
	if s == "" {
		p.addErr(errors.New("cannot parse empty string as an big int"))
		return big.NewInt(0)
	}

	if strings.ContainsAny(s, "eE") {
		f, _, err := new(big.Float).Parse(s, 10)
		if err != nil {
			p.addErr(fmt.Errorf("cannot parse '%s' as an big float", s))
			return big.NewInt(0)
		}
		i, _ := f.Int(nil)
		return i
	}

	i := new(big.Int)
	if _, ok := i.SetString(s, 10); !ok {
		p.addErr(fmt.Errorf("cannot parse '%s' as an big int", s))
		return big.NewInt(0)
	}
	return i
}

func (p *parser) parseAffiliates(s string) []*Affiliate {
	result := make([]*Affiliate, 0)

	if strings.Contains(s, ",") {
		parts := strings.Split(s, ",")
		for _, part := range parts {
			aff := p.parseSingleAffiliate(part)
			if aff == nil {
				return result
			}
			result = append(result, aff)
		}
		return result
	}

	if strings.Contains(s, ":") {
		aff := p.parseSingleAffiliate(s)
		if aff == nil {
			return result
		}
		result = append(result, aff)
		return result
	}

	var nameBuilder strings.Builder
	var bpsBuilder strings.Builder

	flush := func() bool {
		if nameBuilder.Len() == 0 || bpsBuilder.Len() == 0 {
			nameBuilder.Reset()
			bpsBuilder.Reset()
			return true
		}
		name := nameBuilder.String()
		if len(name) > 2 {
			p.addErr(fmt.Errorf("invalid compressed affiliate name '%s'", name))
			return false
		}

		bps, ok := new(big.Int).SetString(bpsBuilder.String(), 10)
		if !ok {
			p.addErr(fmt.Errorf("cannot parse '%s' as an big int", bpsBuilder.String()))
			return false
		}

		result = append(result, &Affiliate{
			Name:       name,
			Bps:        bps,
			Compressed: true,
		})
		nameBuilder.Reset()
		bpsBuilder.Reset()
		return true
	}

	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || r == '-' || r == '_':
			if bpsBuilder.Len() > 0 {
				if !flush() {
					return []*Affiliate{}
				}
			}
			//if nameBuilder.Len() < maxNameLength {
			//	nameBuilder.WriteRune(r)
			//} else {
			//	nameBuilder.Reset()
			//	bpsBuilder.Reset()
			//}
			nameBuilder.WriteRune(r)
		case unicode.IsDigit(r):
			bpsBuilder.WriteRune(r)
		default:
			nameBuilder.Reset()
			bpsBuilder.Reset()
			p.addErr(fmt.Errorf("invalid char '%d'", r))
			return []*Affiliate{}
		}
	}
	if !flush() {
		return []*Affiliate{}
	}

	return result
}

func (p *parser) parseSingleAffiliate(part string) *Affiliate {
	if part == "" {
		return nil
	}
	affiliateParts := strings.SplitN(part, ":", 2)
	if len(affiliateParts) != 2 {
		p.addErr(fmt.Errorf("invalid affiliate '%s'", part))
		return nil
	}
	name := strings.TrimSpace(affiliateParts[0])
	bpsStr := strings.TrimSpace(affiliateParts[1]) // todo negative ?
	if name == "" || bpsStr == "" {
		p.addErr(fmt.Errorf("invalid affiliate '%s'", part))
		return nil
	}
	//if len(name) > maxNameLength {
	//	return nil
	//}
	bps, ok := new(big.Int).SetString(bpsStr, 10)
	if !ok {
		p.addErr(fmt.Errorf("cannot parse '%s' as an big int", bpsStr))
		return nil
	}

	return &Affiliate{
		Name: name,
		Bps:  bps,
	}
}
