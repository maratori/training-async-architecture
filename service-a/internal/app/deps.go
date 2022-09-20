package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/maratori/training-async-architecture/proto-hub/serviceb"
	"github.com/maratori/training-async-architecture/service-a/internal/domain"
	"github.com/maratori/training-async-architecture/service-a/internal/postgres"
)

type ServiceDeps struct {
	queries *postgres.Queries
	client  serviceb.BService
}

func NewServiceDeps(
	queries *postgres.Queries,
	client serviceb.BService,
) *ServiceDeps {
	return &ServiceDeps{
		queries: queries,
		client:  client,
	}
}

func (s *ServiceDeps) SaveX(ctx context.Context, x domain.X) error {
	return s.queries.InsertX(ctx, postgres.X{
		ID:   uuid.MustParse(x.ID),
		Name: x.Name,
	})
}

func (s *ServiceDeps) CallServiceB(ctx context.Context, name string) error {
	_, err := s.client.DoIt(ctx, &serviceb.Request{
		Name: name,
	})
	return err
}
