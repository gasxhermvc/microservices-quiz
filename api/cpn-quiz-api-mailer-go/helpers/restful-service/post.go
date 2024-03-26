package restful

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

func (r *Restful) HttpPostFormDataMultiPart(endpoint string, body *bytes.Buffer, contentType string, jwtToken string) ([]byte, int, error) {
	startProcess := time.Now()
	r.log.Info(r.transId, "Start :: HttpPostFormDataMultiPart")
	r.log.Info(r.transId, "HttpPostFormDataMultiPart.Endpont.Url: ", endpoint)

	//=>กำหนด Request timeout
	timeout := time.Duration(300 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	//=>สร้าง HTTP Method
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		r.log.Error(r.transId, "HttpPostFormDataMultiPart.NewRequest.Error: ", err.Error())
		r.log.Info(r.transId, "End :: HttpPostFormDataMultiPart")
		return nil, 503, err
	}

	//=>กำหนด Content type
	req.Header.Add("Content-Type", contentType)

	//=>ตรวจสอบ JWT
	extractToken := strings.Split(jwtToken, " ")
	if len(extractToken) < 2 {
		r.log.Error(r.transId, "HttpPostFormDataMultiPart.Split.Jwt.Error: ", err.Error())
		r.log.Info(r.transId, "End :: HttpPostFormDataMultiPart")
		return nil, 401, errors.New("token invalid")
	}

	//=>Bind JWT ไปที่ Header
	req.Header.Add("Authorization", html.EscapeString(fmt.Sprintf("%s %s", extractToken[0], extractToken[1])))

	//=>ส่ง Request และตรวจสอบ Request
	resp, err := client.Do(req)
	if err != nil {
		r.log.Error(r.transId, "HttpPostFormDataMultiPart.Client.Do.Error: ", err.Error())
		r.log.Info(r.transId, "End :: HttpPostFormDataMultiPart")
		return nil, 503, err
	}

	defer resp.Body.Close()
	//=>อ่าน Response body
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log.Error(r.transId, "HttpPostFormDataMultiPart.Copy.Response.Error: ", err.Error())
		r.log.Info(r.transId, "End :: HttpPostFormDataMultiPart")
		return bodyText, 503, err
	}

	//=>สกัด Response สำเร็จ
	statusCode := resp.StatusCode
	r.log.Error(r.transId, fmt.Sprintf("HttpPostFormDataMultiPart.Response.Code: %d, Size: %d", statusCode, len(bodyText)), r.log.GetElapsedTime(startProcess))
	r.log.Info(r.transId, "End :: HttpGet")
	return bodyText, statusCode, nil
}
