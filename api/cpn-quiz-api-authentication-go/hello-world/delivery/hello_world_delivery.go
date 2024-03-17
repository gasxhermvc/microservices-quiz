package delivery

import (
	"net/http"
	"time"
	"web-project-template/constant"
	"web-project-template/domain"
	"web-project-template/hello-world/vaildate"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (hello helloWorldDelivery) Hello(c echo.Context) error {
	//=>ต้องใส่ทุกๆ ฟังก์ชัน
	hello.transId = uuid.New().String()
	startProcess := time.Now()
	response := domain.Response{}
	response.TransactionId = hello.transId
	hello.log.Info(hello.transId, "Start Process Hello")

	param := domain.HelloParameter{}

	//=>Check model binder
	if err := c.Bind(&param); err != nil {
		hello.log.Error(hello.transId, err.Error())
		response.Message = constant.BadRequest
		response.Code = constant.BadRequestCode
		return c.JSON(http.StatusBadRequest, response)
	}

	//=>Check model vaildator
	if err := vaildate.CheckRequest(param); err != (domain.Response{}) {
		err.TransactionId = hello.transId
		hello.log.Error(hello.transId, err.Message)
		hello.log.Info(hello.transId, "End Process Hello ", hello.log.GetElapsedTime(startProcess))
		return c.JSON(http.StatusBadRequest, err)
	}

	//=>Bussiness logic
	result := hello.helloWorldUseCase.Hello(param)

	hello.log.Info(hello.transId, "End Process Hello ", hello.log.GetElapsedTime(startProcess))
	return hello.DoResponse(c, result, response)
}
