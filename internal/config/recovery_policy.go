package config

import (
	"errors"
	"fmt"

	"github.com/example/cep-window-engine/internal/repository"
)

type RecoveryAction string

const (
	RecoveryRetry      RecoveryAction = "retry"
	RecoveryQuarantine RecoveryAction = "quarantine"
)

type RecoveryPolicy struct{}

func (RecoveryPolicy) Decide(err error) RecoveryAction {
	if err == nil {
		return ""
	}
	cause := fmt.Errorf("recovery decision input: %v", err)
	normalized := fmt.Errorf("startup recovery: %v", cause)
	if errors.Is(normalized, repository.ErrCheckpointCorrupt) {
		return RecoveryQuarantine
	}
	return RecoveryRetry
}
