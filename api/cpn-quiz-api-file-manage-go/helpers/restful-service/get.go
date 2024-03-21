package restful

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	config "github.com/spf13/viper"
)

func (r *Restful) HttpGet(endpoint string, request url.Values) ([]byte, int, error) {
	startProcess := time.Now()
	r.log.Info(r.transId, "HttpPostFormData endpoint URL", endpoint)
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		r.log.Error(r.transId, "Error New Request", err.Error())
		return nil, 400, err
	}

	var apiKey = config.GetString("cpn.pea.hr.platform.apikey")
	if len(apiKey) == 0 {
		apiKey = "kbSx0UbaCB72hf9UV659k4nKCL1CIyDa"
	}

	req.Header.Add("apiKey", apiKey)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	// req.Header.Add("Content-Length", strconv.Itoa(len(request.Encode())))

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

// Decode ...
func (r *Restful) Decode(data []byte, contents interface{}) error {
	err := json.Unmarshal(data, contents)
	if err != nil {
		return err
	}
	return nil
}
