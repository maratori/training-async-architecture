package app

import (
	"context"

	"github.com/maratori/training-async-architecture/proto-hub/servicea"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DomainService interface {
	DoIt(context.Context, string) error
}

type AService struct {
	svc DomainService
}

func NewAService(svc DomainService) *AService {
	return &AService{
		svc: svc,
	}
}

func (s *AService) DoIt(ctx context.Context, r *servicea.Request) (*emptypb.Empty, error) {
	err := s.svc.DoIt(ctx, r.Name)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
