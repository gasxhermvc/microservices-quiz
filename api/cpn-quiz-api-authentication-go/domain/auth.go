package domain

import "github.com/golang-jwt/jwt/v5"

//=>Repo interfaces.
type AuthRepository interface{}

//=>UseCase interfaces.
type AuthUseCase interface {
	GenerateToken(payload jwt.MapClaims) (response UseCaseResult)
}

//=>App struct.
type UserInfo struct {
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Sub               string `json:"sub"`
}

type Token struct {
	Username   string      `json:"username"`
	UserInfo   *UserInfo   `json:"userInfo"`
	IsEmployee bool        `json:"isEmployee"`
	Permission interface{} `json:"permission"`
	jwt.RegisteredClaims
}
