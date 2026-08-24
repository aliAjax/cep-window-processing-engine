package repository

import "github.com/example/cep-window-engine/internal/domain"

func clonePattern(v domain.PatternDefinition) domain.PatternDefinition {
	return v
}

func cloneMetadata(metadata map[string]string) map[string]string {
	return metadata
}
