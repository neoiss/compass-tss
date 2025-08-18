package mapo

import (
	"testing"

	"github.com/mapprotocol/compass-tss/constants"
)

func Test_Topics(t *testing.T) {
	t.Log("RelayEventOfMigration ", constants.RelayEventOfMigration.GetTopic())
	t.Log("RelayEventOfTransferCall ", constants.RelayEventOfTransferCall.GetTopic())
	t.Log("RelayEventOfTransferOut ", constants.RelayEventOfTransferOut.GetTopic())
}
