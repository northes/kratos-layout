package user

import (
	"github.com/go-kratos/kratos/v2/log"
	userrepo "github.com/northes/kratos-layout/internal/app/repo/user"
	usersvc "github.com/northes/kratos-layout/internal/app/service/user"
	pb "github.com/northes/kratos-layout/internal/app/transport/grpc/api/v1/user"
)

var _ pb.UserServer = (*Handler)(nil)

type Handler struct {
	pb.UnimplementedUserServer
	logger  *log.Helper
	service *usersvc.Service
	repo    userrepo.IRepo
}

func NewHandler(logger log.Logger, service *usersvc.Service, repo userrepo.IRepo) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		service: service,
		repo:    repo,
	}
}
