package restful

import "web-project-template/logger"

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
