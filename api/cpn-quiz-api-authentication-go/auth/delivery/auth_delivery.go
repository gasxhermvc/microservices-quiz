package delivery

import (
	_utils "cpn-quiz-api-authentication-go/utils"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (auth authDelivery) AuthToken(c echo.Context) error {
	tx := uuid.New().String()
	auth.log.Info(tx, "Start :: AuthToken")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	//=>Generate access token
	result := auth.authUsecase.GenerateToken(claims)
	result.Tx = auth.transId
	response := _utils.Response(result)

	if result.Errors != nil {
		//=>Failure
		for _, err := range result.Errors {
			auth.log.Error(tx, fmt.Sprintf("GenerateToken.Error: %s", err))
		}
		auth.log.Error(tx, "End :: AuthToken")
		return c.JSON(response.StatusCode, response)
	}

	//=>Done.
	auth.log.Info(tx, "End :: AuthToken")
	return c.JSON(response.StatusCode, response)
}
