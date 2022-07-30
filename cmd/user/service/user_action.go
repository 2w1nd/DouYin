package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/DouYin/cmd/user/dal/db"
	"github.com/DouYin/cmd/user/pack"
	"github.com/DouYin/kitex_gen/user"
	"github.com/DouYin/pkg/constants"
	"github.com/DouYin/pkg/errno"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
)

type UserActionService struct {
	ctx context.Context
}

func NewUserActionService(ctx context.Context) *UserActionService {
	return &UserActionService{
		ctx: ctx,
	}
}

func (c *UserActionService) Create(req *user.CreateUserRequest) (*user.User, error) {
	_, err := db.QueryUser(c.ctx, req.Username)
	if err == nil {
		return nil, errno.UserErr.WithMsg("user is exist")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	h := sha256.New()
	if _, err = io.WriteString(h, req.Password+constants.UserSalt); err != nil {
		return nil, err
	}
	password := fmt.Sprintf("%x", h.Sum(nil))

	usr, err := db.Create(c.ctx, &db.User{
		UserId:   int64(uuid.New().ID()),
		Name:     req.Username,
		Username: req.Username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return pack.User(usr), nil
}
