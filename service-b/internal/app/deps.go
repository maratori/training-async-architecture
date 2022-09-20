package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/maratori/training-async-architecture/proto-hub/servicea"
	"github.com/maratori/training-async-architecture/service-b/internal/domain"
	"github.com/maratori/training-async-architecture/service-b/internal/postgres"
)

type ServiceDeps struct {
	queries *postgres.Queries
	client  servicea.AService
}

func NewServiceDeps(
	queries *postgres.Queries,
	client servicea.AService,
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

func (s *ServiceDeps) CallServiceA(ctx context.Context, name string) error {
	_, err := s.client.DoIt(ctx, &servicea.Request{
		Name: name,
	})
	return err
}
