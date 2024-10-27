package auth

import (
	"errors"
	"order_api/config"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	jwtService *JWTService
}

func NewAuthService(config *config.Config) *AuthService {
	return &AuthService{
		jwtService: NewJWTService(config),
	}
}

// Login 处理用户登录
func (s *AuthService) Login(username, password string) (string, error) {
	// TODO: 实现实际的用户验证逻辑
	// 这里仅作示例，实际应用中需要查询数据库验证用户
	if username == "admin" && password == "admin123" {
		return s.jwtService.GenerateToken("admin-user-id", "admin")
	}
	return "", ErrInvalidCredentials
}

// ValidateToken 验证令牌
func (s *AuthService) ValidateToken(token string) (*Claims, error) {
	return s.jwtService.ValidateToken(token)
}
