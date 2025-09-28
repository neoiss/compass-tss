package mapo

import (
	"testing"

	"github.com/mapprotocol/compass-tss/constants"
)

func Test_Topics(t *testing.T) {
	t.Log("EventOfBridgeIn ", constants.EventOfBridgeIn.GetTopic())
	t.Log("EventOfBridgeOut ", constants.EventOfBridgeOut.GetTopic())
}
