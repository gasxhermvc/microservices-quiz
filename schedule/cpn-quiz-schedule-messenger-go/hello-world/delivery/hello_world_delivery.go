package delivery

// func (h *helloWorldDelivery) TriggerJob(c echo.Context) error {
// 	// token := c.Get("user").(*jwt.Token).Claims.(*domain.Token)

// 	// h.helloWorldUsecase.SetUsername(token.Username)
// 	// h.helloWorldUsecase.SetTransactionID(uuid.New().String())

// 	// h.log.Info(h.helloWorldUsecase.GetTransactionID(), "Start Trigger Back Office SCS Update Hello World")
// 	// isRunning, err := h.helloWorldUsecase.CheckBackOfficeIsRunning()

// 	// if err != nil {
// 	// 	h.log.Error(h.helloWorldUsecase.GetTransactionID(), "Back Office Erros : "+err.Error())
// 	// 	return c.JSON(http.StatusInternalServerError, domain.Response{
// 	// 		TransactionId: h.helloWorldUsecase.GetTransactionID(),
// 	// 		Code:          constant.InternalServerErrorCode,
// 	// 		Message:       constant.InternalServerError,
// 	// 	})
// 	// }

// 	// if isRunning {
// 	// 	h.log.Error(h.helloWorldUsecase.GetTransactionID(), "Can't Trigger Because Back Office is Running ")
// 	// 	return c.JSON(http.StatusBadRequest, domain.Response{
// 	// 		TransactionId: h.helloWorldUsecase.GetTransactionID(),
// 	// 		Code:          constant.BadRequestCode,
// 	// 		Message:       "Can't Trigger Because Back Office is Running ",
// 	// 	})
// 	// }

// 	// go h.helloWorldUsecase.CronJob()
// 	h.log.Info(h.helloWorldUsecase.GetTransactionID(), "Finish Trigger Back Office SCS Update Hello World")
// 	return c.JSON(http.StatusOK, domain.Response{

// 		TransactionId: h.helloWorldUsecase.GetTransactionID(),
// 		Code:          constant.SuccessCode,
// 		Message:       constant.Success,
// 	})
// }
