// internal/middleware/jwt.go
package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const JWTSessionConst = "__JWT_SESSION__"

var (
	// 这个密钥在实际应用中应该从配置文件或环境变量中读取
	secretKey = []byte("abcd")
)

// 自定义JWT声明结构体
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type JSession struct {
	Uid       int    `json:"uid"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint) (string, error) {
	// 设置 JWT 声明
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token 有效期 24 小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并获得完整编码后的字符串 Token
	return token.SignedString(secretKey)
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证 Token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		// 解析Token
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 将用户ID存储在上下文中，以便后续处理器使用
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
