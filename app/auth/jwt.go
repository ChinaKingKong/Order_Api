// Package auth 提供JWT认证相关的功能实现
package auth

import (
	"errors"
	"order_api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义JWT相关错误
var (
	ErrInvalidToken = errors.New("无效的令牌")
	ErrExpiredToken = errors.New("令牌已过期")
)

// Claims 定义JWT的声明结构
type Claims struct {
	UserID string `json:"user_id"` // 用户ID
	Role   string `json:"role"`    // 用户角色
	jwt.RegisteredClaims
}

// JWTService JWT服务结构体
type JWTService struct {
	config *config.Config // 配置信息
}

// NewJWTService 创建新的JWT服务实例
func NewJWTService(config *config.Config) *JWTService {
	return &JWTService{
		config: config,
	}
}

// GenerateToken 生成JWT令牌
// userID: 用户ID
// role: 用户角色
func (s *JWTService) GenerateToken(userID, role string) (string, error) {
	// 创建JWT声明
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.config.JWT.TokenExpiryHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建并签名令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.SecretKey))
}

// ValidateToken 验证JWT令牌
// tokenString: JWT令牌字符串
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌并返回声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
