package main

import (
	"context"
	"github.com/DouYin/cmd/user/pack"
	"github.com/DouYin/cmd/user/service"
	"github.com/DouYin/kitex_gen/user"
	"github.com/DouYin/pkg/errno"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// CreateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	resp = new(user.CreateUserResponse)
	// TODO: validate the params...

	usr, err := service.NewUserActionService(ctx).Create(req)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	resp.UserId = usr.Id
	return resp, nil
}

// MGetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) MGetUser(ctx context.Context, req *user.MGetUserRequest) (resp *user.MGetUserResponse, err error) {
	resp = new(user.MGetUserResponse)
	users, err := service.NewUserInfoService(ctx).MGet(req)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	resp.Users = users
	return resp, nil
}

// CheckUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CheckUser(ctx context.Context, req *user.CheckUserRequest) (resp *user.CheckUserResponse, err error) {
	resp = new(user.CheckUserResponse)

	uid, err := service.NewUserInfoService(ctx).Check(req)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	resp.UserId = uid
	return resp, nil
}

// IsExist implements the UserServiceImpl interface.
func (s *UserServiceImpl) IsExist(ctx context.Context, req *user.IsExistByIdRequest) (resp *user.IsExistByIdResponse, err error) {
	// TODO: Your code here...
	return
}
