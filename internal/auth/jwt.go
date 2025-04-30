package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// https://jwt.io/introduction
type JWTAuthenticator struct {
	secret string
	aud    string // 接收者
	iss    string // 签发者
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret: secret, iss: iss, aud: aud}
}

func (a *JWTAuthenticator) GenerateToken(claim jwt.Claims) (string, error) {
	// 生成一个token，使用SigningMethodHS256签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// 使用密钥对 Token 进行签名，生成最终的 JWT 字符串
	// a.secret 是签名的秘密密钥，必须保密
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	// 1. 解析 JWT 字符串。
	//    jwt.Parse 是解析和验证的核心函数。
	//    它需要 JWT 字符串、一个用于获取验证密钥的回调函数，以及可选的验证选项。
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(a.secret), nil
	},
		// 2. 添加额外的验证选项 (Options)。
		//    这些选项会根据 JWT 中的标准声明进行验证。
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.aud),
		// 2.4. 验证签名算法。
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
