package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type BService struct{}

func NewBService() *BService {
	return &BService{}
}

func (*BService) DoIt(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
