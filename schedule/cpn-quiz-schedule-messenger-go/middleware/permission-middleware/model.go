package permissions

import (
	"cpn-quiz-schedule-messenger-go/logger"
	"sync"
)

type PermissionMiddleware struct {
	log   *logger.PatternLogger
	Mutex sync.Mutex
}

type MatchRoute struct {
	Route    string
	Resource string
}
