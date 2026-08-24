package partition

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
	assignments := make(map[int]string, len(partitions))
	for i, partition := range partitions {
		assignments[partition] = active[i%len(active)].ID()
	}
	return assignments, nil
}
