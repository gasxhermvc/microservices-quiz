package restful

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func (r *Restful) HttpGetSummary(endpoint string, body *bytes.Buffer) ([]byte, int, error) {
	startProcess := time.Now()
	r.log.Info(r.transId, "HttpPostFormData endpoint URL", endpoint)
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	//endpoint = "http://172.30.200.184/CPN/CJ20N?wbs=C64FSUNCM0192021"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		r.log.Error(r.transId, "Error New Request", err.Error())
		return nil, 400, err
	}

	resp, err := client.Do(req)
	if err != nil {
		r.log.Error(r.transId, "Error Client Do", err.Error())
		return nil, 400, err
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log.Error(r.transId, "Error ReadAll Body Text", err.Error())
		return nil, 400, err
	}

	statusCode := resp.StatusCode
	r.log.Info(r.transId, "HttpPostFormData Response :", r.log.GetElapsedTime(startProcess))
	return bodyText, statusCode, nil
}

func (r *Restful) HttpGetSummaryUnLog(endpoint string, body *bytes.Buffer) ([]byte, int, error) {
	//startProcess := time.Now()
	//r.log.Info(r.transId, "HttpPostFormData endpoint URL", endpoint)
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	//endpoint = "http://172.30.200.184/CPN/CJ20N?wbs=C64FSUNCM0192021"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		//r.log.Error(r.transId, "Error New Request", err.Error())
		return nil, 400, err
	}

	/* var apiKey = config.GetString("cpn.pea.hr.platform.apikey")
	if len(apiKey) == 0 {
		apiKey = "kbSx0UbaCB72hf9UV659k4nKCL1CIyDa"
	}

	req.Header.Add("apiKey", apiKey) */
	// req.Header.Add("Accesskey", accessKey)
	// req.Header.Add("Appkey", appKey)
	//req.Header.Add("username", username)
	//req.Header.Add("Content-Type", "application/json; charset=utf-8")
	// req.Header.Add("Content-Length", strconv.Itoa(len(request.Encode())))

	resp, err := client.Do(req)
	if err != nil {
		//r.log.Error(r.transId, "Error Client Do", err.Error())
		return nil, 400, err
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		//r.log.Error(r.transId, "Error ReadAll Body Text", err.Error())
		return nil, 400, err
	}

	statusCode := resp.StatusCode
	//r.log.Info(r.transId, "HttpPostFormData Response :", r.log.GetElapsedTime(startProcess))
	return bodyText, statusCode, nil
}

func (r *Restful) HttpGet(endpoint string, request GetRequest) ([]byte, int, error) {
	startProcess := time.Now()
	r.log.Info(r.transId, "HttpPostFormData endpoint URL", endpoint)
	// Request Timeout
	timeout := time.Duration(180 * time.Minute)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		r.log.Error(r.transId, "Error New Request", err.Error())
		return nil, 400, err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if request.Headers.ApiKey != nil {
		req.Header.Add("api-key", *request.Headers.ApiKey)
	}

	if request.Headers.Authorization != nil {
		req.Header.Add("Authorization", *request.Headers.Authorization)
	}

	if request.Headers.XClientID != nil {
		req.Header.Add("x-client-id", *request.Headers.XClientID)
	}

	resp, err := client.Do(req)
	if err != nil {
		r.log.Error(r.transId, "Error Client Do", err.Error())
		return nil, 400, err
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log.Error(r.transId, "Error ReadAll Body Text", err.Error())
		return nil, 400, err
	}

	statusCode := resp.StatusCode

	r.log.Info(r.transId, "HttpGetFormData Response :", r.log.GetElapsedTime(startProcess))
	return bodyText, statusCode, nil
}

// Decode ...
func (r *Restful) Decode(data []byte, contents interface{}) error {
	err := json.Unmarshal(data, contents)
	if err != nil {
		return err
	}
	return nil
}
