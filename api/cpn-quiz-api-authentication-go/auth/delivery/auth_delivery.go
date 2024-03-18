package delivery

import (
	"cpn-quiz-api-authentication-go/utils"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nqmt/goerror"
)

func (auth authDelivery) AuthToken(c echo.Context) error {
	tx := uuid.New().String()
	auth.log.Info(tx, "Start :: AuthToken")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	//=>Generate access token
	result := auth.authUsecase.GenerateToken(claims)
	if result.Error != nil {
		//=>Failure
		auth.log.Error(tx, fmt.Sprintf("GenerateToken.Error: %s", result.Error.Error()))
		auth.log.Error(tx, "End :: AuthToken")
		return goerror.EchoErrorReturn(result.Error, c, tx)
	}

	//=>Done.
	result.Tx = tx
	auth.log.Info(tx, "End :: AuthToken")
	return c.JSON(http.StatusOK, utils.Response(result))
}
