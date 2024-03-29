package task

import (
	"cpn-quiz-schedule-messenger-go/constant"
	"cpn-quiz-schedule-messenger-go/database"
	"cpn-quiz-schedule-messenger-go/domain"
	"cpn-quiz-schedule-messenger-go/logger"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron"
	config "github.com/spf13/viper"
)

type helloWorldTask struct {
	emailMessengerUseCase domain.EmailMessengerUseCase
	log                   *logger.PatternLogger
	rdbCon                *database.RedisDatabase
	transId               string
	cron                  *cron.Cron
}

func NewHelloWorldHandler(emailMessengerUseCase domain.EmailMessengerUseCase, rdbConn *database.RedisDatabase, log *logger.PatternLogger, cron *cron.Cron) {
	handler := &helloWorldTask{
		emailMessengerUseCase: emailMessengerUseCase,
		rdbCon:                rdbConn,
		log:                   log,
		cron:                  cron,
	}

	handler.Run()
}

func (task *helloWorldTask) Run() {
	guid := uuid.New().String()
	startApp := time.Now()
	startAppProcess := startApp
	startAppTimeFormat := startApp.Format(config.GetString("cpn.quiz.format.shortdate"))
	task.log.Info(guid, "Start :: Task.Run")

	var counter struct {
		success int64
		error   int64
		total   int64
	}

	var timeSchedule = "*/5 * * * *"

	running := false
	_ = task.cron.AddFunc(timeSchedule, func() {
		if running {
			// task.log.Info(task.transId, "Job exit, If another process running.")
			return
		}

		running = true
		for {
			task.transId = uuid.New().String()
			currentTime := time.Now()
			startProcess := currentTime
			timeFormat := currentTime.Format(config.GetString("cpn.quiz.format.shortdate"))

			task.log.Info(task.transId, "Start Dequeue to Send Email ", fmt.Sprintf("at %s", timeFormat))

			dequeue, deQueueErr := task.rdbCon.Dequeue(config.GetString("cpn.quiz.queue.email"))
			if deQueueErr != nil {
				if deQueueErr.Error() == "queue is empty" {
					time.Sleep(time.Second * 5)
					continue
				}

				if dequeue != "" {
					task.log.Info(task.transId, "For test ===> Retry to queue...")
					//=>Error แบบยัง Dequeue สำเร็จ, จะต้องนำรายการกลับเข้า Queue
					task.rdbCon.Enqueue(config.GetString("cpn.quiz.queue.email"), dequeue)
					task.log.Info(task.transId, "For test ===> Done...")
				}
				counter.error++
				counter.total = counter.success + counter.error
				task.log.Info(task.transId, "Task.Schedule.Dequeue.Error: ", deQueueErr.Error())
				task.log.Info(task.transId, "End Dequeue to Send Email ", fmt.Sprintf("at %s use timne %d", timeFormat, task.log.GetElapsedTime(startProcess)))
				time.Sleep(time.Second * 5)
				continue
			}

			if dequeue != "" {
				task.log.Info(task.transId, fmt.Sprintf("For Test ===> Process Queue: %s", dequeue))
			}

			var provider domain.EmailProvider
			var parameter domain.EmailQueueParameter
			err := json.Unmarshal([]byte(dequeue), &parameter)
			if err != nil {
				if dequeue != "" {
					task.log.Info(task.transId, "For test ===> Retry to queue...")
					//=>Error แบบยัง Dequeue สำเร็จ, จะต้องนำรายการกลับเข้า Queue
					task.rdbCon.Enqueue(config.GetString("cpn.quiz.queue.email"), dequeue)
					task.log.Info(task.transId, "For test ===> Done...")
				}
				counter.error++
				counter.total = counter.success + counter.error
				task.log.Info(task.transId, "Task.Schedule.Deserialize.Error: ", err.Error())
				task.log.Info(task.transId, "End Dequeue to Send Email ", fmt.Sprintf("at %s use timne %d", timeFormat, task.log.GetElapsedTime(startProcess)))
				time.Sleep(time.Second * 5)
				continue
			}

			//=>Setup parameter after dequeue.
			provider.Config = *buildEmailConfiguration()
			provider.Parameter = parameter
			result := task.emailMessengerUseCase.SendEmail(&provider)
			result.Tx = task.transId
			if result.Errors != nil {
				//=>สำหรับ Customize สถานะว่าส่ง Email ไม่สำเร็จจะเอารายการกลับเข้าคิว
				if result.StatusCode == constant.ServiceUnavailableCode {
					if dequeue != "" {
						task.log.Info(task.transId, "For test ===> Retry to queue...")
						//=>Error แบบยัง Dequeue สำเร็จ, จะต้องนำรายการกลับเข้า Queue
						task.rdbCon.Enqueue(config.GetString("cpn.quiz.queue.email"), dequeue)
						task.log.Info(task.transId, "For test ===> Done...")
					}
				}

				//=>Loop errors
				for _, err := range result.Errors {
					task.log.Info(task.transId, "Task.Schedule.SendEmail.Error: ", err)
				}

				counter.error++
				counter.total = counter.success + counter.error
				task.log.Info(task.transId, "End Dequeue to Send Email ", fmt.Sprintf("at %s use timne %d", timeFormat, task.log.GetElapsedTime(startProcess)))
				time.Sleep(time.Second * 5)
				continue
			}
			//=>Send completed
			counter.success++
			counter.total = counter.success + counter.error
			task.log.Info(task.transId, fmt.Sprintf("Counter: Success: %d, Error: %d, Total: %d", counter.success, counter.error, counter.total))
			task.log.Info(task.transId, "End Dequeue to Send Email ", fmt.Sprintf("at %s use timne %d", timeFormat, task.log.GetElapsedTime(startProcess)))
			time.Sleep(5 * time.Second)
		}
		running = false
	})
	task.cron.Start()
	task.log.Info(guid, "Stop Task Run :: ", fmt.Sprintf("at %s use timne %d", startAppTimeFormat, task.log.GetElapsedTime(startAppProcess)))
}

func buildEmailConfiguration() *domain.EmailConfig {
	emailConfig := new(domain.EmailConfig)
	emailConfig.Server = config.GetString("cpn.quiz.api.mailer.email.server")
	emailConfig.Port = config.GetInt("cpn.quiz.api.mailer.email.port")
	emailConfig.EnableSSL = config.GetBool("cpn.quiz.api.mailer.email.enablessl")
	emailConfig.SendWithCredential = config.GetBool("cpn.quiz.api.mailer.email.sendwithcredential")
	emailConfig.Skip = config.GetBool("cpn.quiz.api.mailer.email.skip")
	emailConfig.DefaultCredential = config.GetBool("cpn.quiz.api.mailer.email.defaultcredential")
	emailConfig.Username = config.GetString("cpn.quiz.api.mailer.email.username")
	emailConfig.Password = config.GetString("cpn.quiz.api.mailer.email.password")
	emailConfig.SenderAddress = config.GetString("cpn.quiz.api.mailer.email.senderaddress")
	emailConfig.FromParameter = config.GetString("cpn.quiz.api.mailer.email.fromparameter")
	emailConfig.ToParameter = config.GetString("cpn.quiz.api.mailer.email.toparameter")
	emailConfig.CcParameter = config.GetString("cpn.quiz.api.mailer.email.ccparameter")
	emailConfig.BccParameter = config.GetString("cpn.quiz.api.mailer.email.bccparameter")
	emailConfig.SubjectParameter = config.GetString("cpn.quiz.api.mailer.email.subjectparameter")
	emailConfig.BodyParameter = config.GetString("cpn.quiz.api.mailer.email.bodyparameter")
	emailConfig.PriorityParameter = config.GetString("cpn.quiz.api.mailer.email.priorityparameter")
	emailConfig.AttachmentFileParameter = config.GetString("cpn.quiz.api.mailer.email.attachmentparameter")

	return emailConfig
}
