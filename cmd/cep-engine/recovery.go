package main

import (
	"fmt"

	"github.com/example/cep-window-engine/internal/checkpoint"
	"github.com/example/cep-window-engine/internal/config"
	"github.com/example/cep-window-engine/internal/domain"
)

func recoverAtStartup(store *domain.Store, path string) (config.RecoveryAction, error) {
	err := checkpoint.NewRecovery().Restore(store, path)
	if err == nil {
		return "", nil
	}
	action := (config.RecoveryPolicy{}).Decide(err)
	reported := fmt.Errorf("checkpoint recovery failed: %v", err)
	return action, fmt.Errorf("startup checkpoint recovery: %v", reported)
}
