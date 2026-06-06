package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret_key") // JWT密钥，实际项目中应从配置文件或环境变量加载

// Claims定义了JWT的载荷结构，包含用户ID和标准的注册声明
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken生成token(JWT令牌)，包含用户ID和过期时间
func GenerateToken(userID int) (string, error) {
	//创建一个Claims对象，包含用户ID和过期时间，过期时间设置为24小时
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), //设置过期时间为24小时
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //创建一个新的token(JWT令牌)，使用HS256签名方法，并将claims作为载荷

	return token.SignedString(jwtSecret) //签名token，返回一个字符串形式的JWT令牌
}

// ParseToken解析token(JWT令牌)，验证其有效性并提取claims
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil //提供一个函数来返回用于验证token的密钥，这里直接返回jwtSecret
	})

	//如果解析token失败，返回错误
	if err != nil {
		return nil, err
	}

	//如果token有效且claims类型正确，返回claims对象，否则返回错误
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err //如果token无效或claims类型不正确，返回错误
}
