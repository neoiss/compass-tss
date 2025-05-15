package thorchain

import (
	"fmt"

	"github.com/mapprotocol/compass-tss/common"
)

type DonateMemo struct{ MemoBase }

func (m DonateMemo) String() string {
	return fmt.Sprintf("DONATE:%s", m.Asset)
}

func (p *parser) ParseDonateMemo() (DonateMemo, error) {
	asset := p.getAsset(1, true, common.EmptyAsset)
	return DonateMemo{
		MemoBase: MemoBase{TxType: TxDonate, Asset: asset},
	}, p.Error()
}
