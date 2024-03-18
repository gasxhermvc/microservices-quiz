package task

import (
	"cpn-quiz-schedule-messenger-go/domain"
	"cpn-quiz-schedule-messenger-go/logger"

	"github.com/robfig/cron"
)

type helloWorldTask struct {
	helloWorldUseCase domain.HelloWorldUseCase
	log               *logger.PatternLogger
	transId           string
	cron              *cron.Cron
}

func NewHelloWorldHandler(helloWorldUseCase domain.HelloWorldUseCase, log *logger.PatternLogger, cron *cron.Cron) {
	handler := &helloWorldTask{
		helloWorldUseCase: helloWorldUseCase,
		log:               log,
		cron:              cron,
	}

	handler.Run()
}

func (task *helloWorldTask) Run() {
	// 0 1 * * *  --> Cron job every day at 1am is a commonly used cron schedule.
	// */10 * * * *  --> every 10 seconds

	// var timeSchedule = "*/10 * * * *"

	// _ = task.cron.AddFunc(timeSchedule, func() {
	// 	task.transId = uuid.New().String()
	// 	currentTime := time.Now()
	// 	timeFormat := currentTime.Format("2006-01-02")
	// 	startProcess := currentTime

	// 	task.log.Info(task.transId, "Start Hello ", fmt.Sprintf("at %s", timeFormat))
	// 	err := task.helloWorldUseCase.CronJob()
	// 	if err != nil {
	// 		task.log.Error(task.transId, "Error Process CronJob ", err.Error())
	// 	}

	// 	task.log.Info(task.transId, "End Hello ", fmt.Sprintf("at %s use timne %d", timeFormat, task.log.GetElapsedTime(startProcess)))
	// })

	// task.cron.Start()
}
