package service

import (
	"context"
	"errors"
	"time"

	"micro-service-platform/services/user-service/internal/model"
	"micro-service-platform/services/user-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *model.UserCreateRequest) (*model.User, error)
	Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserLoginResponse, error)
	GetUser(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, id uint, req *model.UserUpdateRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, req *model.UserListRequest) (*model.UserListResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*model.User, error)
}

type userService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: "your-secret-key", // 实际项目中应该从配置文件读取
	}
}

func (s *userService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.User, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Status:   model.UserStatusActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	// 根据用户名获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 检查用户状态
	if user.Status == model.UserStatusBanned {
		return nil, errors.New("user is banned")
	}
	if user.Status == model.UserStatusInactive {
		return nil, errors.New("user is inactive")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 生成JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.UserLoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *userService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id uint, req *model.UserUpdateRequest) (*model.User, error) {
	// 获取现有用户
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, req *model.UserListRequest) (*model.UserListResponse, error) {
	return s.userRepo.List(ctx, req)
}

func (s *userService) ValidateToken(ctx context.Context, tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return s.userRepo.GetByID(ctx, userID)
	}

	return nil, errors.New("invalid token")
}

func (s *userService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
