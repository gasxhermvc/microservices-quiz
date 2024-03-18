package restful

import (
	"cpn-quiz-schedule-messenger-go/logger"
	"net/url"
)

type Restful struct {
	log     *logger.PatternLogger
	transId string
}

type GetRequest struct {
	Headers GetRequestHeaders
	Params  url.Values
}

type GetRequestHeaders struct {
	ApiKey        *string `json:"api-key"`
	Authorization *string `json:"Authorization"`
	XClientID     *string `json:"x-client-id"`
}

func NewRestful(log *logger.PatternLogger, transId string) *Restful {
	rest := &Restful{
		log:     log,
		transId: transId,
	}

	return rest
}
