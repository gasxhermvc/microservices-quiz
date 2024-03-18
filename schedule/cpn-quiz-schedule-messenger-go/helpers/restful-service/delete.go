package restful

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func (r *Restful) HttpDeleteFormDataMultiPart(endpoint string, authJwt, clientId, apiKey string) ([]byte, int, error) {
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		r.log.Error(r.transId, "Error New Request ", err.Error())
		return nil, 400, err
	}
	fmt.Println("x-client-id :", clientId)
	fmt.Println("x-api-key :", apiKey)
	fmt.Println("Authorization :", authJwt)

	if authJwt != "" {
		req.Header.Add("Authorization", authJwt)
	}
	if clientId != "" {
		req.Header.Add("x-client-id", clientId)
	}
	if apiKey != "" {
		req.Header.Add("x-api-key", apiKey)
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
		return bodyText, 400, err
	}

	statusCode := resp.StatusCode

	r.log.Info(r.transId, "HttpPostFormData Response :", string(bodyText))
	return bodyText, statusCode, nil
}