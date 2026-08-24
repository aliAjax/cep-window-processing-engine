package replayload

import "context"

type SegmentFetcher interface {
	Fetch(context.Context, string) ([]ReplayRecord, error)
}

type ReplayService struct {
	repository SegmentFetcher
}

func NewService(repository SegmentFetcher) *ReplayService {
	return &ReplayService{repository: repository}
}

func (s *ReplayService) Run(ctx context.Context, request ReplayRequest) ([]ReplayRecord, error) {
	if s == nil || s.repository == nil {
		return nil, nil
	}
	return s.repository.Fetch(context.Background(), request.Segment)
}
