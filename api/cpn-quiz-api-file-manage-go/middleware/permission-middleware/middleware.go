package permissions

import (
	"cpn-quiz-api-file-manage-go/constant"
	"cpn-quiz-api-file-manage-go/domain"
	"cpn-quiz-api-file-manage-go/logger"
	"net/http"
	"slices"

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
				claims := token.Claims.(jwt.MapClaims)
				permissions, found := claims["realm_access"].(map[string]interface{})
				if !found {
					permissions = claims["permission"].(map[string]interface{})
				}
				roles := permissions["roles"].([]interface{})

				var convertRoles []string
				for _, role := range roles {
					convertRoles = append(convertRoles, role.(string))
				}

				isMatches := false
				for _, match := range matches {
					if slices.Contains(convertRoles, match) {
						isMatches = true
						break
					}
				}

				if !isMatches {
					msg = "Access denied., permission not matches."
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
