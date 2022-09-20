package app

import (
	"context"

	"github.com/maratori/training-async-architecture/proto-hub/servicea"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AService struct{}

func NewAService() *AService {
	return &AService{}
}

func (*AService) DoIt(context.Context, *servicea.Request) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
