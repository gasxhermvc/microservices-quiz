package usecase

import (
	"cpn-quiz-api-authentication-go/constant"
	"cpn-quiz-api-authentication-go/domain"
	"cpn-quiz-api-authentication-go/logger"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/spf13/viper"
)

// =>Setup time for web
var webTimeExpiredCpnQuiz = 3 * time.Hour

type authUseCase struct {
	authRepository domain.AuthRepository
	log            *logger.PatternLogger
}

func NewAuthUseCase(authRepository domain.AuthRepository, log *logger.PatternLogger) domain.AuthUseCase {
	return &authUseCase{
		authRepository: authRepository,
		log:            log,
	}
}

func (auth *authUseCase) GenerateToken(payload jwt.MapClaims) domain.UseCaseResult {
	result := domain.UseCaseResult{}

	//=>Setup life time.
	tokenExpireDuration := webTimeExpiredCpnQuiz
	if os.Args[0] != "prod" {
		tokenExpireDuration = webTimeExpiredCpnQuiz * 8
	}

	//=>Create user info.
	userInfo := domain.UserInfo{
		PreferredUsername: payload["preferred_username"].(string),
		Email:             payload["email"].(string),
		GivenName:         payload["given_name"].(string),
		FamilyName:        payload["family_name"].(string),
	}

	//=>Create expired time.
	expireAt := jwt.NumericDate{
		Time: time.Now().AddDate(0, 0, 0).Add(tokenExpireDuration),
	}

	//=>Create claim.
	claims := &domain.Token{
		Username:         payload["preferred_username"].(string),
		UserInfo:         &userInfo,
		Permission:       payload["realm_access"],
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &expireAt},
		Sub:              payload["sub"].(string),
		Aud:              "cpn-quiz",
		Iat:              time.Now(),
		Iss:              payload["iss"].(string),
	}

	//=>Generate token by secretkey in db config.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(config.GetString("cpn.quiz.api.jwt.secretkey")))

	if err != nil {
		//=>Generate error.
		result.Errors = append(result.Errors, err.Error())
		result.StatusCode = constant.UnAuthorizationCode
		result.Message = constant.UnAuthorization
		result.Success = false
		return result
	}

	//=>Done.
	result.Result = accessToken
	result.Success = true
	result.Message = constant.Success
	result.StatusCode = constant.SuccessCode

	return result
}
