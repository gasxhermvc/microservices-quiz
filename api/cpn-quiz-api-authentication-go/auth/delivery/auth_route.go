package delivery

import (
	"context"
	"cpn-quiz-api-authentication-go/domain"
	"cpn-quiz-api-authentication-go/logger"
	"errors"
	"fmt"

	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	jwk "github.com/lestrrat-go/jwx/jwk"
	config "github.com/spf13/viper"
)

func getKey(token *jwt.Token) (interface{}, error) {
	endpoint := config.GetString("cpm.quiz.sso.endpoint")
	keySet, err := jwk.Fetch(context.Background(), endpoint+"/realms/cpn-quiz/protocol/openid-connect/certs")
	if err != nil {
		return nil, err
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have a key ID in the kid field")
	}

	key, found := keySet.LookupKeyID(keyID)
	if !found {
		return nil, fmt.Errorf("unable to find key %q", keyID)
	}

	var pubkey interface{}
	if err := key.Raw(&pubkey); err != nil {
		return nil, fmt.Errorf("unable to get the public key. Error: %s", err.Error())
	}

	return pubkey, nil
}

type authDelivery struct {
	authUsecase domain.AuthUseCase
	log         *logger.PatternLogger
	transId     string
}

func NewAuthDelivery(e *echo.Echo, authUsecase domain.AuthUseCase, log *logger.PatternLogger) {
	handler := &authDelivery{
		authUsecase: authUsecase,
		log:         log,
	}

	r := e.Group(config.GetString("service.endpoint"))

	g := r.Group("/auth")
	eConfig := echojwt.Config{
		KeyFunc: getKey,
	}

	g.Use(echojwt.WithConfig(eConfig))
	g.POST("/token", handler.AuthToken)
}
