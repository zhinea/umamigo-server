package libjwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/utils"
)

var SigningMethod = jwt.SigningMethodHS256

var SignatureKey = []byte(utils.Cfg.AppSecret)

// CreateToken creates a new JWT token
func CreateToken(payload entity.SessionClaims) string {

	claims := entity.JWTSessionClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "UmamiGO",
		},
		SessionClaims: payload,
	}

	token := jwt.NewWithClaims(SigningMethod, claims)
	ss, err := token.SignedString(SignatureKey)

	if err != nil {
		panic(err)
	}

	return ss
}

// ParseToken parses a JWT token
func ParseToken(tokenString string) (*entity.SessionClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.JWTSessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SignatureKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*entity.JWTSessionClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &claims.SessionClaims, nil
}
