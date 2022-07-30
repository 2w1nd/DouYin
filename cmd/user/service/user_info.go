package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/DouYin/cmd/user/dal/db"
	"github.com/DouYin/kitex_gen/user"
	"github.com/DouYin/pkg/constants"
	"github.com/DouYin/pkg/errno"
	"github.com/bytedance/gopkg/util/logger"
	"io"
)

type UserInfoService struct {
	ctx context.Context
}

func NewUserInfoService(ctx context.Context) *UserInfoService {
	return &UserInfoService{
		ctx: ctx,
	}
}

func (c *UserInfoService) Check(req *user.CheckUserRequest) (uid int64, err error) {
	logger.Info("check user")
	h := sha256.New()
	if _, err = io.WriteString(h, req.Password+constants.UserSalt); err != nil {
		return 0, err
	}
	password := fmt.Sprintf("%x", h.Sum(nil))
	username := req.Username
	u, err := db.QueryUser(c.ctx, username)
	if err != nil {
		return 0, err
	}

	if u.Password != password {
		return 0, errno.UserErr.WithMsg("username or password is wrong")
	}
	return u.UserId, nil
}

func (c *UserInfoService) MGet(req *user.MGetUserRequest) ([]*user.User, error) {
	//if len(req.TargetUserIds) == 0 {
	//	return make([]*user.User, 0), nil
	//}
	//urs, err := db.MGet(m.ctx, req.TargetUserIds)
	//if err != nil {
	//	return nil, err
	//}
	//if len(urs) == 0 {
	//	return nil, errno.UserErr.WithMsg("user isn't exist")
	//}
	//users := pack.Users(urs)
	//for i, u := range users {
	//	countInfo, err := rpc.RelationInfo(m.ctx, &relation.InfoRequest{
	//		UserId:       req.UserId,
	//		TargetUserId: u.Id,
	//	})
	//	if err != nil {
	//		return users, err
	//	}
	//	users[i].FollowCount = countInfo.FollowCount
	//	users[i].FollowerCount = countInfo.FollowerCount
	//	users[i].IsFollow = countInfo.IsFollow
	//}
	return nil, nil
}

func (c *UserInfoService) IsExist(req *user.IsExistByIdRequest) (bool, error) {
	return db.IsExist(c.ctx, req.UserId)
}
