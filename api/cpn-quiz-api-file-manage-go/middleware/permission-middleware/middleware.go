package permissions

import (
	"cpn-quiz-api-file-manage-go/constant"
	"cpn-quiz-api-file-manage-go/domain"
	"cpn-quiz-api-file-manage-go/logger"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewPermissionMiddleware(log *logger.PatternLogger) *PermissionMiddleware {
	return &PermissionMiddleware{
		log: log,
	}
}

func (perm *PermissionMiddleware) AuthorizePermissions(matches ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			perm.Mutex.Lock()
			defer perm.Mutex.Unlock()

			msg := ""
			transId := uuid.New().String()

			token, found := c.Get("user").(*jwt.Token)
			if found {

				x1, _ := token.Claims.GetExpirationTime()
				fmt.Println(x1)
				x2, _ := token.Claims.GetIssuedAt()
				fmt.Println(x2)

				x3, _ := token.Claims.GetAudience()
				fmt.Println(x3)

				x4, _ := token.Claims.GetIssuer()
				fmt.Println(x4)

				x5, _ := token.Claims.GetSubject()
				fmt.Println(x5)

				x6, _ := token.Claims.GetNotBefore()
				fmt.Println(x6)
				// claims := token.Claims.(jwt.MapClaims)
				// claims["permissions"]

				// permissions := claims["permission"].(map[string]interface{})
				// // if !found {
				// // 	permissions = claims["permission"].(map[string]interface{})
				// // }
				// roles := permissions["roles"].([]interface{})

				// var convertRoles []string
				// for _, role := range roles {
				// 	convertRoles = append(convertRoles, role.(string))
				// }

				// isMatches := false
				// for _, match := range matches {
				// 	if slices.Contains(convertRoles, match) {
				// 		isMatches = true
				// 		break
				// 	}
				// }

				// if !isMatches {
				// 	msg = "Access denied., permission not matches."
				// 	perm.log.Error(transId, msg)
				// 	return c.JSON(http.StatusForbidden, domain.Response{
				// 		TransactionId: transId,
				// 		Message:       constant.AccessDenied,
				// 		Code:          constant.AccessDeniedCode,
				// 		ErrorResponse: domain.ErrorResponse{
				// 			Error: []string{msg},
				// 		},
				// 	})
				// }
			} else {
				clientId := c.Request().Header.Get("x-client-id")
				authorization := c.Request().Header.Get("x-api-key")

				if clientId == "" && authorization == "" {
					msg = "Un authorization, invalid or expired jwt."
					perm.log.Error(transId, msg)
					return c.JSON(http.StatusForbidden, domain.Response{
						TransactionId: transId,
						Message:       constant.AccessDenied,
						Code:          constant.AccessDeniedCode,
						ErrorResponse: domain.ErrorResponse{
							Error: []string{msg},
						},
					})
				}
			}

			//=>Next or Error
			if err := next(c); err != nil {
				return c.JSON(http.StatusUnauthorized, domain.Response{
					TransactionId: transId,
					Message:       constant.UnAuthorization,
					Code:          constant.UnAuthorizationCode,
					ErrorResponse: domain.ErrorResponse{
						Error: []string{err.Error()},
					},
				})
			}

			return nil
		}
	}
}
