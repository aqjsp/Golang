package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"micro-service-platform/services/user-service/internal/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *model.UserListRequest) (*model.UserListResponse, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redisClient *redis.Client) UserRepository {
	// 自动迁移数据库表（仅在数据库连接可用时）
	if db != nil {
		db.AutoMigrate(&model.User{})
	}

	return &userRepository{
		db:    db,
		redis: redisClient,
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	// 清除相关缓存
	r.clearUserCache(user.ID)
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	// 先从缓存获取
	cacheKey := r.getUserCacheKey(id)
	cachedUser, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user model.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// 从数据库获取
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	// 缓存到Redis
	r.cacheUser(&user)

	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}

	// 清除缓存
	r.clearUserCache(user.ID)
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return err
	}

	// 清除缓存
	r.clearUserCache(id)
	return nil
}

func (r *userRepository) List(ctx context.Context, req *model.UserListRequest) (*model.UserListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	query := r.db.WithContext(ctx).Model(&model.User{})

	// 添加搜索条件
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	var users []model.User
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}

	return &model.UserListResponse{
		Users: users,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// 缓存相关方法
func (r *userRepository) getUserCacheKey(id uint) string {
	return fmt.Sprintf("user:%d", id)
}

func (r *userRepository) cacheUser(user *model.User) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return
	}

	cacheKey := r.getUserCacheKey(user.ID)
	r.redis.Set(context.Background(), cacheKey, userJSON, 30*time.Minute)
}

func (r *userRepository) clearUserCache(id uint) {
	cacheKey := r.getUserCacheKey(id)
	r.redis.Del(context.Background(), cacheKey)
}
