package permissions

import (
	"cpn-quiz-api-file-manage-go/logger"
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
