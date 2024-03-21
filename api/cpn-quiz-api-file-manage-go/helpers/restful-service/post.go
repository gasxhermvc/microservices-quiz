package restful

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	config "github.com/spf13/viper"
)

func (r *Restful) HttpPost(endpoint string, request url.Values) ([]byte, int, error) {

	r.log.Info(r.transId, "HttpPostFormData endpoint URL", endpoint)
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(request.Encode()))
	if err != nil {
		r.log.Error(r.transId, "Error New Request", err.Error())
		return nil, 400, err
	}

	var apiKey = config.GetString("cpn.hr.platform.apikey")
	req.Header.Add("apiKey", apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(request.Encode())))

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

	r.log.Info(r.transId, "HttpPostFormData Response :", string(bodyText))
	return bodyText, statusCode, nil
}
