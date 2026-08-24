package recovery

import "fmt"

type Consumer struct {
	Coordinator *Coordinator
}

func (c Consumer) Commit(epochID string) error {
	if c.Coordinator == nil {
		return fmt.Errorf("recovery coordinator is unavailable")
	}
	epoch, ok := c.Coordinator.Get(epochID)
	if !ok {
		return fmt.Errorf("recovery epoch not found")
	}
	if !epoch.Complete() || !epoch.CanTransition(StateCommitted) {
		return fmt.Errorf("recovery epoch is not ready to commit")
	}
	epoch.State = StateCommitted
	c.Coordinator.Put(epoch)
	return nil
}
