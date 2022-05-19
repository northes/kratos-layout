package user

import (
	"context"

	"github.com/northes/kratos-layout/internal/app/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	errorx "github.com/northes/kratos-layout/internal/app/pkg/errors"
	"github.com/northes/kratos-layout/internal/app/pkg/validator"
)

type CreateRequest struct {
	Name  string `json:"name"`
	Age   uint8  `json:"age"`
	Phone string `json:"phone"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("名称不能为空")),
		validation.Field(&r.Phone, validation.By(validator.IsMobilePhone)),
	)
}

type CreateResponse struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Age   uint8  `json:"age"`
	Phone string `json:"phone"`
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorx.New(errorx.ValidateErrorCode, err.Error())
	}

	m := &model.User{
		Name:  req.Name,
		Age:   req.Age,
		Phone: req.Phone,
	}

	if _, err := s.repo.Create(ctx, m); err != nil {
		s.logger.Errorf("%s: %s", model.ErrDataStoreFailed, err)
		return nil, errorx.New(
			errorx.ServerErrorCode,
			model.ErrDataStoreFailed.Error(),
		)
	}

	resp := &CreateResponse{
		Id:    m.ID,
		Name:  m.Name,
		Age:   m.Age,
		Phone: m.Phone,
	}

	return resp, nil
}
