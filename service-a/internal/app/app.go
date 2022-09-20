package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type AService struct{}

func NewAService() *AService {
	return &AService{}
}

func (*AService) DoIt(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
