package repository

import (
	"github.com/DouYin/common/constant"
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"gorm.io/gorm"
)

type FollowRepository struct {
	Base BaseRepository
	User UserRepository
}

func (r *FollowRepository) AddFollow(follow model.Follow) bool {
	where := model.Follow{
		FollowedUserId: follow.FollowedUserId,
	}

	if _, err := r.User.GetFirstUser(where); err != nil {
		return false
	}

	if err := r.Base.Create(&follow); err != nil {
		return false
	}
	return true
}

func (r *FollowRepository) UpdateFollowUserId(where interface{}, out interface{}) bool {
	db := global.DB.Where(where)
	if err := db.Model(out).Where(where).Update("is_deleted", 0).Error; err != nil {
		return false
	}
	return true
}

func (r *FollowRepository) DeleteFollowUserId(where interface{}) bool {
	var follow model.Follow
	if err := r.Base.DeleteSoftByID(where, &follow); err != nil {
		return false
	}
	return true
}

func (r *FollowRepository) GetFollowerListByUserId(id uint64) ([]model.Follow, error) {
	var followList []model.Follow
	query := global.DB.
		Model(model.Follow{}).
		Where("followed_user_id = ?", id).
		Preload("User")
	err := query.Debug().Find(&followList).Error
	if err != nil {
		return followList, err
	}
	return followList, nil
}

func (r *FollowRepository) GetFollowedOrFollowUserWithUserId(id uint64, Type int) ([]dto.FollowDto, error) {
	var followedList []dto.FollowDto

	var subQuery, subQuery1 *gorm.DB
	var sqlQ string
	subQuery = global.DB.Model(model.Follow{}).Select("user_id", "followed_user_id", "is_deleted")
	subQuery1 = global.DB.Model(model.Follow{}).Select("user_id", "followed_user_id", "is_deleted")
	if Type == constant.Followed {
		subQuery = subQuery.Where("followed_user_id = ?", id)
		subQuery1 = subQuery1.Where("user_id = ?", id)
		sqlQ = "LEFT JOIN douyin_user on a.user_id = douyin_user.user_id"
	} else if Type == constant.Follow {
		subQuery = subQuery.Where("user_id = ?", id)
		subQuery1 = subQuery1.Where("followed_user_id = ?", id)
		sqlQ = "LEFT JOIN douyin_user on a.followed_user_id = douyin_user.user_id"
	}

	err := global.DB.Debug().
		Select("a.user_id ",
			"username as name",
			"a.followed_user_id ",
			"CASE WHEN a.is_deleted = 0 THEN true ELSE false END as followed_a ",
			"CASE WHEN b.is_deleted = 0 THEN true ELSE false END as followed_b ",
			"follow_count ",
			"follower_count ").
		Where("a.is_deleted = 0").
		Table("(?) as a "+
			"LEFT JOIN (?) as b "+
			"ON a.user_id = b.followed_user_id AND a.followed_user_id = b.user_id "+sqlQ, subQuery, subQuery1).
		Find(&followedList).Error

	if err != nil {
		return followedList, err
	}
	return followedList, nil
}
