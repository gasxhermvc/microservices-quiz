package restful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (r *Restful) HttpGet(endpoint string, request url.Values, filename string, token string) ([]byte, int, error) {
	startProcess := time.Now()
	r.log.Info(r.transId, "Start :: HttpGet")
	r.log.Info(r.transId, "HttpGet.Endpont.Url: ", endpoint)

	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		r.log.Error(r.transId, "HttpGet.NewRequest.Error: ", err.Error())
		r.log.Info(r.transId, "End :: HttpGet")
		return nil, 503, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		r.log.Error(r.transId, "HttpGet.Client.Do.Error: ", err.Error())
		r.log.Info(r.transId, "End :: HttpGet")
		return nil, 503, err
	}

	defer resp.Body.Close()
	//=>จัดการ Temp
	temp := "./temp"
	stats, err := os.Stat(temp)
	if os.IsNotExist(err) {
		os.Mkdir(temp, os.FileMode(0755))
	}
	r.log.Info(r.transId, "Create temp: "+stats.Name())

	fullAccess := "./" + strings.Replace(strings.Replace(strings.Replace(filepath.Join(temp, filename), "\\\\", "//", -1), "\\", "/", -1), "//", "/", -1)
	file, err := os.Create(fullAccess)
	if err != nil {
		r.log.Error(r.transId, "HttpGet.Create.File.Error: ", err.Error())
	}
	size, err := io.Copy(file, resp.Body)
	if err != nil {
		r.log.Error(r.transId, "HttpGet.Body.Copy.Error: ", err.Error())
		r.log.Info(r.transId, "End :: HttpGet")
		return nil, 503, err
	}
	file.Close()

	src, err := os.Open(fullAccess)
	if err != nil {
		r.log.Error(r.transId, "HttpGet.Open.File.Error: ", err.Error())
		r.log.Info(r.transId, "End :: HttpGet")
		return nil, 503, err
	}
	buf := bytes.NewBuffer(nil)
	fsize, err := io.Copy(buf, src)
	if err != nil {
		r.log.Error(r.transId, "HttpGet.Copy.Error: ", err.Error())
		r.log.Info(r.transId, "End :: HttpGet")
		return nil, 503, err
	}
	defer (func() {
		src.Close()

		err = os.Remove(fullAccess)
		if err != nil {
			r.log.Error(r.transId, "HttpGet.Delete.Error: ", err.Error())
		}
	})()

	r.log.Info(r.transId, "HttpGet.Copy.Size: ", fsize)
	statusCode := resp.StatusCode
	r.log.Error(r.transId, fmt.Sprintf("HttpGet.Response.Code: %d, Size: %d", statusCode, size), r.log.GetElapsedTime(startProcess))
	r.log.Info(r.transId, "End :: HttpGet")
	return buf.Bytes(), statusCode, nil
}

// Decode ...
func (r *Restful) Decode(data []byte, contents interface{}) error {
	err := json.Unmarshal(data, contents)
	if err != nil {
		return err
	}
	return nil
}
