package app

import (
	"context"

	"github.com/maratori/training-async-architecture/proto-hub/serviceb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DomainService interface {
	DoIt(context.Context, string) error
}

type BService struct {
	svc DomainService
}

func NewBService(svc DomainService) *BService {
	return &BService{
		svc: svc,
	}
}

func (s *BService) DoIt(ctx context.Context, r *serviceb.Request) (*emptypb.Empty, error) {
	err := s.svc.DoIt(ctx, r.Name)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
