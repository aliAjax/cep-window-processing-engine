package backfill

import "context"

type Fetcher interface {
	Fetch(context.Context, string) ([]Record, error)
}

type Service struct {
	repository Fetcher
}

func NewService(repository Fetcher) *Service {
	return &Service{repository: repository}
}

func (s *Service) Run(ctx context.Context, request Request) ([]Record, error) {
	if s == nil || s.repository == nil {
		return nil, nil
	}
	return s.repository.Fetch(context.Background(), request.Shard)
}
