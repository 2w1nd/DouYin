package repository

import (
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
)

type FollowRepository struct {
	Base BaseRepository
}

func (r *FollowRepository) AddFollow(follow model.Follow) bool {
	if err := r.Base.Create(&follow); err != nil {
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

func (r *FollowRepository) GetFollowedUserWithUserId(id uint64) ([]dto.FollowDto, error) {
	var followedList []dto.FollowDto

	subQuery := global.DB.
		Model(model.Follow{}).
		Select("user_id", "followed_user_id", "is_deleted").
		Where("followed_user_id = ?", id)
	subQuery1 := global.DB.
		Model(model.Follow{}).
		Select("user_id", "followed_user_id", "is_deleted").
		Where("user_id = ?", id)
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
			"ON a.user_id = b.followed_user_id AND a.followed_user_id = b.user_id "+
			"LEFT JOIN douyin_user on a.user_id = douyin_user.user_id", subQuery, subQuery1).
		Find(&followedList).Error
	if err != nil {
		return followedList, err
	}
	return followedList, nil
}

func (r *FollowRepository) GetFollowUserWithUserId(id uint64) ([]dto.FollowDto, error) {
	var followList []dto.FollowDto

	subQuery1 := global.DB.
		Model(model.Follow{}).
		Select("user_id", "followed_user_id", "is_deleted").
		Where("followed_user_id = ?", id)
	subQuery := global.DB.
		Model(model.Follow{}).
		Select("user_id", "followed_user_id", "is_deleted").
		Where("user_id = ?", id)
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
			"ON a.user_id = b.followed_user_id AND a.followed_user_id = b.user_id "+
			"LEFT JOIN douyin_user on a.user_id = douyin_user.user_id", subQuery, subQuery1).
		Find(&followList).Error
	if err != nil {
		return followList, err
	}
	return followList, nil
}
