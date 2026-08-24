package partition

import "fmt"

type Worker struct {
	Registry *Registry
}

func (w Worker) Apply(assignments map[int]string) error {
	if w.Registry == nil {
		return fmt.Errorf("partition registry is unavailable")
	}
	for partition, owner := range assignments {
		w.Registry.Assign(partition, owner)
	}
	return nil
}
