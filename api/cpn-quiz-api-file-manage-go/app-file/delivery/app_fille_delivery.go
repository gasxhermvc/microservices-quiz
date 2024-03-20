package delivery

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (appFile appFileDelivery) UploadFile(c echo.Context) error {
	appFile.transId = uuid.New().String()
	return c.JSON(http.StatusOK, "OK")
}

func (appFile appFileDelivery) RemoveFile(c echo.Context) error {
	appFile.transId = uuid.New().String()
	return c.JSON(http.StatusOK, "OK")
}

func (appFile appFileDelivery) PreviewFile(c echo.Context) error {
	appFile.transId = uuid.New().String()
	return c.JSON(http.StatusOK, "OK")
}

func (appFile appFileDelivery) DownloadFile(c echo.Context) error {
	appFile.transId = uuid.New().String()
	return c.JSON(http.StatusOK, "OK")
}
