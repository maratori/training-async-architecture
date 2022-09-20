package app

import (
	"context"

	"github.com/maratori/training-async-architecture/proto-hub/serviceb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BService struct{}

func NewBService() *BService {
	return &BService{}
}

func (*BService) DoIt(context.Context, *serviceb.Request) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
