package usecase

import (
	"cpn-quiz-schedule-messenger-go/domain"
	"cpn-quiz-schedule-messenger-go/logger"
)

type helloWorldUseCase struct {
	helloWorldRepository domain.HelloWorldRespository
	log                  *logger.PatternLogger
	runDate              string
	transId              string
	username             string
	PID                  int
}

func NewHelloWorldUseCase(helloWorldRepository domain.HelloWorldRespository, log *logger.PatternLogger) domain.HelloWorldUseCase {
	return &helloWorldUseCase{
		helloWorldRepository: helloWorldRepository,
		log:                  log,
	}
}

// func (hello helloWorldUseCase) CronJob() error {
// 	// fmt.Println("trigg")
// 	hello.log.Info(hello.transId, "Start Hello World...")
// 	localePrinter := message.NewPrinter(language.English)

// 	//=>Setup
// 	hello.runDate = time.Now().Format("2006-01-02 15:04:05")
// 	pId, _ := hello.InsertLog()
// 	hello.PID = pId
// 	result := domain.UpdateResult{}

// 	hello.InsertLogDetail(pId, "เริ่มการประมวลผล Hello World")
// 	hello.UpdateLog(pId, constant.LOG_INPROGRESS, result.Total, result.Success, result.Error)

// 	hello.InsertLogDetail(pId, "เริ่มการบันทึกข้อมูลลง Database")
// 	processStart := time.Now()
// 	//=>something
// 	hello.InsertLogDetail(pId, "สิ้นสุดการบันทึกข้อมูลลง Database")
// 	timeElapsed := hello.log.GetElapsedTimeSec(processStart)
// 	hello.InsertLogDetail(pId, localePrinter.Sprintf("สถิติการบันทึกข้อมูลลง Database ใช้เวลาบันทึกข้อมูลทั้งสิ้น : %.2f วินาที", timeElapsed))
// 	hello.callSummarizeNotification(result)
// 	hello.log.Info(hello.transId, "Hello World...")
// 	hello.InsertLogDetail(pId, "สิ้นสุดการประมวลผล Hello World")
// 	hello.UpdateLog(pId, constant.LOG_SUCESS, result.Total, result.Success, result.Error)
// 	hello.SetUsername("")
// 	return nil
// }

// func (hello helloWorldUseCase) displayUpdateStatus(result domain.UpdateResult) {
// 	p := message.NewPrinter(language.English)
// 	hello.log.Info(hello.transId, fmt.Sprintf(`Summarize...
// 		Total: %s
// 		Success: %s
// 		Error: %s`, p.Sprintf("%d", result.Total), p.Sprintf("%d", result.Success), p.Sprintf("%d", result.Error)))
// }

// func (hello helloWorldUseCase) callSummarizeNotification(result domain.UpdateResult) {
// 	p := message.NewPrinter(language.English)
// 	lineToken := config.GetString("cpn.api.notification.line.token")
// 	lineNotification := notification.NewLineNotify(lineToken)
// 	params := make(map[string]interface{})
// 	msgDetail := fmt.Sprintf(`
// 	ผลการบันทึกลง DB
// 	📓 Total: %s
// 	✅ Success: %s
// 	❌ Error: %s`, p.Sprintf("%d", result.Total), p.Sprintf("%d", result.Success), p.Sprintf("%d", result.Error))
// 	params = notification.LineNotiMsgTemplate(`Back Office - Hello World`, true, msgDetail)
// 	lineNotification.Create(params)
// 	err1 := lineNotification.Send()
// 	if err1 != nil {
// 		hello.log.Error(hello.transId, err1.Error())
// 	}
// }

// func (hello helloWorldUseCase) callErroNotification(errMsg string) {
// 	lineToken := config.GetString("cpn.api.notification.line.token")
// 	lineNotification := notification.NewLineNotify(lineToken)
// 	params := make(map[string]interface{})
// 	msgDetail := fmt.Sprintf(`
// 	%s`, errMsg)
// 	params = notification.LineNotiMsgTemplate(`Back Office - Hello World`, false, msgDetail)
// 	lineNotification.Create(params)
// 	err1 := lineNotification.Send()
// 	if err1 != nil {
// 		hello.log.Error(hello.transId, err1.Error())
// 	}
// }

// func (hello helloWorldUseCase) InsertLogDetail(logStatusId int, detail string) error {
// 	params := make(map[string]interface{})
// 	// params := make(map[string]string)
// 	params["LOG_PROGRAM_STATUS_ID"] = logStatusId
// 	params["DETAIL"] = detail

// 	// hello.log.Debug(hello.transId, detail)
// 	err := hello.helloWorldRepository.InsertLogDetail(params)
// 	if err != nil {
// 		hello.log.Error(hello.transId, err.Error())
// 		return err
// 	}
// 	return nil
// }

// func (hello helloWorldUseCase) InsertLog() (int, error) {
// 	params := make(map[string]interface{})

// 	params["PROGRAM_ID"] = config.GetInt("cpn.application.programId")
// 	params["RUN_DATE"] = hello.runDate
// 	if len(hello.username) != 0 {
// 		params["USERNAME"] = hello.username
// 	}

// 	pID, err := hello.helloWorldRepository.InsertLog(params)
// 	if err != nil {
// 		return pID, err
// 	}

// 	return pID, nil
// }

// func (hello helloWorldUseCase) UpdateLog(pId int, status int, total int, success int, fail int) error {
// 	params := make(map[string]interface{})
// 	params["ID"] = pId
// 	params["STATUS_ID"] = status
// 	params["TOTAL"] = total
// 	params["TOTAL_SUCCESS"] = success
// 	params["TOTAL_ERROR"] = fail

// 	if len(hello.username) != 0 {
// 		params["USERNAME"] = hello.username
// 	}

// 	err := hello.helloWorldRepository.UpdateLog(params)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (hello helloWorldUseCase) CheckBackOfficeIsRunning() (bool, error) {
// 	logProgramStatus, err := hello.helloWorldRepository.GetLogProgramStatus()
// 	if err != nil {
// 		return false, err
// 	}

// 	if len(logProgramStatus) == 0 {
// 		return false, errors.New("Not Found")
// 	}

// 	if logProgramStatus[0].STATUS_ID == 3 {
// 		return true, nil
// 	}

// 	return false, nil
// }

// func (hello *helloWorldUseCase) GetTransactionID() string {
// 	return hello.transId
// }

// func (hello *helloWorldUseCase) SetTransactionID(transactionID string) {
// 	hello.transId = transactionID

// }

// func (hello *helloWorldUseCase) SetUsername(username string) {
// 	hello.username = username
// }
