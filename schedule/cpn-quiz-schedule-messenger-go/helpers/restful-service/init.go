package restful

import "cpn-quiz-schedule-messenger-go/logger"

type Restful struct {
	log     *logger.PatternLogger
	transId string
}

func NewRestful(log *logger.PatternLogger, transId string) *Restful {
	rest := &Restful{
		log:     log,
		transId: transId,
	}

	return rest
}
