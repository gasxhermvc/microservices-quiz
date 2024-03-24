package permissions

import (
	"cpn-quiz-api-authentication-go/logger"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

type (
	PermissionMiddleware struct {
		log   *logger.PatternLogger
		Mutex sync.Mutex
	}

	PayloadData struct {
		Username   string      `json:"username"`
		UserInfo   *UserInfo   `json:"userInfo"`
		Permission interface{} `json:"permission"`
		jwt.RegisteredClaims
	}

	UserInfo struct {
		PreferredUsername string `json:"preferred_username"`
		Email             string `json:"email"`
		GivenName         string `json:"given_name"`
		FamilyName        string `json:"family_name"`
		Sub               string `json:"sub"`
	}

	Permissions struct {
		Roles []string `json:"roles"`
	}
)
