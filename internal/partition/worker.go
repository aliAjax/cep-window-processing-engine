package partition

import "fmt"

type Worker struct {
	Registry *Registry
}

func (w Worker) Apply(assignments map[int]string) error {
	for partition, owner := range assignments {
		if owner == "" {
			return fmt.Errorf("partition %d has empty owner", partition)
		}
		w.Registry.Assign(partition, owner)
	}
	return nil
}
