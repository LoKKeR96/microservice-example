package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/palantir/stacktrace"
)

type UserAuthenticationService interface {
	Authenticate(username string, password string) (string, error)
}

type userAuthenticationService struct {
	SecretAuthKey string
}

func NewUserAuthenticationService(secretAuthKey string) UserAuthenticationService {
	return &userAuthenticationService{
		SecretAuthKey: secretAuthKey,
	}
}

const (
	USERNAME        = "member1"
	PASSWORD        = "password123"
	TokenExpiration = 17*time.Hour + 30*time.Minute
)

// just a simple helper function to generate a jwt token with single user in code
func (uas *userAuthenticationService) Authenticate(username string, password string) (string, error) {
	if username != USERNAME || password != PASSWORD {
		return "", stacktrace.Propagate(
			fmt.Errorf("error on validating credentials"),
			"error on validating credentials",
		)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(TokenExpiration).Unix(),
	})

	return claims.SignedString([]byte(uas.SecretAuthKey))
}
