package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/northes/kratos-layout/internal/app/model"
	"gorm.io/gorm"
	"time"
)

var _ IRepo = (*Repository)(nil)

type IRepo interface {
	FindOneById(ctx context.Context, id uint64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Save(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, user *model.User) error
}

type Repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

var (
	cacheKeyFormat = model.User{}.TableName() + "_%d"
	cacheExpire    = 3600
)

func (r *Repository) FindOneById(ctx context.Context, id uint64) (*model.User, error) {
	m := new(model.User)

	cacheValue, err := r.rdb.Get(
		ctx,
		fmt.Sprintf(cacheKeyFormat, id),
	).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, err
		}
	}

	if cacheValue != nil {
		if err = jsoniter.Unmarshal(cacheValue, m); err != nil {
			return nil, err
		}
		return m, nil
	}

	err = r.db.Where("id = ?", id).First(m).Error
	if err != nil {
		return nil, err
	}

	cacheValue, err = jsoniter.Marshal(m)
	if err != nil {
		return nil, err
	}

	err = r.rdb.Set(
		ctx,
		fmt.Sprintf(cacheKeyFormat, id),
		string(cacheValue),
		time.Duration(cacheExpire)*time.Second,
	).Err()
	if err != nil {
		return nil, err
	}

	return m, nil

}

func (r *Repository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	cacheValue, err := jsoniter.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = r.rdb.Set(
		ctx,
		fmt.Sprintf(cacheKeyFormat, user.ID),
		string(cacheValue),
		time.Duration(cacheExpire)*time.Second,
	).Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}

	cacheValue, err := jsoniter.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = r.rdb.Set(
		ctx,
		fmt.Sprintf(cacheKeyFormat, user.ID),
		string(cacheValue),
		time.Duration(cacheExpire)*time.Second,
	).Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Delete(ctx context.Context, user *model.User) error {
	if err := r.db.Delete(user).Error; err != nil {
		return err
	}

	err := r.rdb.Del(
		ctx,
		fmt.Sprintf(cacheKeyFormat, user.ID),
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB, rdb *redis.Client) *Repository {
	return &Repository{
		db:  db,
		rdb: rdb,
	}
}
