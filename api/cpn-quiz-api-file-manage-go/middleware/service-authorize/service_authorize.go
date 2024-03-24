package service_authorize

import (
	"cpn-quiz-api-file-manage-go/constant"
	"cpn-quiz-api-file-manage-go/domain"
	"cpn-quiz-api-file-manage-go/logger"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewCustomAuthorizeGuard(log *logger.PatternLogger) *CustomAuthorizeGuard {
	return &CustomAuthorizeGuard{
		Log: log,
	}
}

func NewCustomAuthorizeGuardOnly(only string, log *logger.PatternLogger) *CustomAuthorizeGuard {
	return &CustomAuthorizeGuard{
		Only: only,
		Log:  log,
	}
}

func (cag *CustomAuthorizeGuard) AuthorizeGuard() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			msg := ""
			transId := uuid.New().String()
			cag.Mutex.Lock()
			defer cag.Mutex.Unlock()

			token := c.Request().Header.Get("Authorization")
			if token == "" {
				clientId := c.Request().Header.Get("x-client-id")
				authorization := c.Request().Header.Get("x-api-key")

				//=>Check client is not empty
				if clientId == "" {
					msg = "Un authorization, required is 'x-client-id'"
					cag.Log.Error(transId, msg)
					return c.JSON(http.StatusUnauthorized, domain.Response{
						TransactionId: transId,
						Message:       msg,
						Code:          constant.UnAuthorizationCode,
					})
				}

				//=>Check apiKey is not empty
				if authorization == "" {
					msg = "Un authorization, required is 'x-api-key'"
					cag.Log.Error(transId, msg)
					return c.JSON(http.StatusUnauthorized, domain.Response{
						TransactionId: transId,
						Message:       constant.UnAuthorization,
						Code:          constant.UnAuthorizationCode,
						ErrorResponse: domain.ErrorResponse{
							Error: []string{msg},
						},
					})
				}

				//=>Check is required 'Only' on route
				if cag.Only != "" && !strings.EqualFold(strings.ToLower(clientId), strings.ToLower(cag.Only)) {
					msg = "Un authorization, Permission denied."
					cag.Log.Error(transId, msg)
					return c.JSON(http.StatusUnauthorized, domain.Response{
						TransactionId: transId,
						Message:       constant.UnAuthorization,
						Code:          constant.UnAuthorizationCode,
						ErrorResponse: domain.ErrorResponse{
							Error: []string{msg},
						},
					})
				}

				apiKey := authorization
				//=>ตรวจสอบ x-api-key
				if len(apiKey) < 1 {
					msg = "Un authorization, please check your ApiKey"
					cag.Log.Error(transId, msg)
					return c.JSON(http.StatusUnauthorized, domain.Response{
						TransactionId: transId,
						Message:       constant.UnAuthorization,
						Code:          constant.UnAuthorizationCode,
						ErrorResponse: domain.ErrorResponse{
							Error: []string{msg},
						},
					})
				}

				//=>Get from global list data
				secretKey := AuthorizeApiKeyLists
				match := new(SecretKey)

				for _, v := range secretKey {
					if strings.EqualFold(strings.ToLower(v.ClientId), strings.ToLower(clientId)) &&
						strings.EqualFold(v.ApiKey, apiKey) {
						match.ClientId = v.ClientId
						match.ApiKey = v.ApiKey
					}
				}

				//=>Check match is empty
				if match.ClientId == "" || match.ApiKey == "" {
					msg = "Un authorization, client-id or api-key incorrect."
					cag.Log.Error(transId, msg)
					return c.JSON(http.StatusUnauthorized, domain.Response{
						TransactionId: transId,
						Message:       msg,
						Code:          constant.UnAuthorizationCode,
					})
				}

				//=>Check api not matching
				if match.ApiKey != apiKey {
					msg = "Un authorization, ClientId or ApiKey incorrect."
					cag.Log.Error(transId, msg)
					return c.JSON(http.StatusUnauthorized, domain.Response{
						TransactionId: transId,
						Message:       msg,
						Code:          constant.UnAuthorizationCode,
					})
				}

				//=>Next or Error
				if err := next(c); err != nil {
					msg = err.Error()
					cag.Log.Error(transId, msg)
					return c.JSON(http.StatusUnauthorized, domain.Response{
						TransactionId: transId,
						Message:       msg,
						Code:          constant.UnAuthorizationCode,
					})
				}
			}
			return nil
		}
	}
}
