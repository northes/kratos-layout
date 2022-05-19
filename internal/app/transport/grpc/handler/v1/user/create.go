package user

import (
	"context"

	"github.com/northes/kratos-layout/internal/app/service/user"
	pb "github.com/northes/kratos-layout/internal/app/transport/grpc/api/v1/user"
)

func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	svcReq := user.CreateRequest{
		Name:  req.Name,
		Age:   uint8(req.Age),
		Phone: req.Phone,
	}

	_, err := h.service.Create(ctx, svcReq)
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{}, nil
}
