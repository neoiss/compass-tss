package runners

import (
	"sync"
	"time"

	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"

	"github.com/mapprotocol/compass-tss/common"
)

// SolvencyCheckProvider methods that a SolvencyChecker implementation should have
type SolvencyCheckProvider interface {
	GetHeight() (int64, error)
	ShouldReportSolvency(height int64) bool
	ReportSolvency(height int64) error
}

// SolvencyCheckRunner when a chain get marked as insolvent , and then get halt automatically , the chain client will stop scanning blocks , as a result , solvency checker will
// not report current solvency status to relay anymore, this method is to ensure that the chain client will continue to do solvency check even when the chain has been halted
func SolvencyCheckRunner(chain common.Chain,
	provider SolvencyCheckProvider,
	bridge shareTypes.Bridge,
	stopper <-chan struct{},
	wg *sync.WaitGroup,
	backOffDuration time.Duration,
) {

}
