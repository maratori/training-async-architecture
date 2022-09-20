package domain

import (
	"context"

	"github.com/google/uuid"
)

type ServiceDeps interface {
	SaveX(context.Context, X) error
	CallServiceB(context.Context, string) error
}

type Service struct {
	deps ServiceDeps
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		deps: deps,
	}
}

func (s *Service) DoIt(ctx context.Context, name string) error {
	x := X{
		ID:   uuid.NewString(),
		Name: name,
	}

	err := s.deps.SaveX(ctx, x)
	if err != nil {
		return err
	}

	if name == "call-b" {
		errC := s.deps.CallServiceB(ctx, name)
		if errC != nil {
			return errC
		}
	}

	return nil
}
