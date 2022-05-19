package user

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/northes/kratos-layout/internal/app/repo/user"
)

type IService interface {
	Create(ctx context.Context, request CreateRequest) (*CreateResponse, error)
	// TODO Update,Delete,Detail,List
}

var _ IService = (*Service)(nil)

type Service struct {
	logger *log.Helper
	repo   user.IRepo
}

func NewService(logger log.Logger, repo user.IRepo) *Service {
	return &Service{
		logger: log.NewHelper(logger),
		repo:   repo,
	}
}
