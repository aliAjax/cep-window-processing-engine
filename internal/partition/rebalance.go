package partition

import "fmt"

type Owner interface {
	ID() string
}

type Rebalancer struct{}

func (Rebalancer) Plan(partitions []int, owners []Owner) (map[int]string, error) {
	active := make([]Owner, 0, len(owners))
	for _, owner := range owners {
		if owner != nil && owner.ID() != "" {
			active = append(active, owner)
		}
	}
	if len(active) == 0 && len(partitions) > 0 {
		return nil, fmt.Errorf("no partition owners available")
	}
	assignments := make(map[int]string, len(partitions))
	for i, partition := range partitions {
		assignments[partition] = active[i%len(active)].ID()
	}
	return assignments, nil
}
