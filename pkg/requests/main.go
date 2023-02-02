package requests

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type HttpRequestI interface {
	Request(method string, path string, header string, body string, json []byte, token string) (resp []byte, statusCode int, err error)
}

type httpClient struct {
	baseUrl string
	timeOut int
}

func NewHttpClient(baseUrl string, timeOut int) HttpRequestI {
	return &httpClient{
		baseUrl: baseUrl,
		timeOut: timeOut,
	}
}

func (h *httpClient) Request(method string, path string, header string, body string, json []byte, token string) (resp []byte, statusCode int, err error) {
	var (
		res *http.Response
	)
	client := &http.Client{Timeout: time.Duration(h.timeOut) * time.Second}

	var req *http.Request

	switch header {
	case "post/params":
		req, err = http.NewRequest(method, path, nil)
		if err != nil {
			return nil, 500, err
		}

		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		req.URL.RawQuery = body

	case "application/x-www-form-urlencoded":
		req, err = http.NewRequest(method, path, strings.NewReader(body))
		if err != nil {
			return nil, 500, err
		}
		req.Header.Add("Content-Length", strconv.Itoa(len(body)))
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		req.Header.Add("Content-Type", header)

	case "application/json":
		req, err = http.NewRequest(method, path, bytes.NewBuffer(json))
		if err != nil {
			return nil, 500, err
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		req.Header.Add("Content-Type", header)

	default:
		return nil, 500, errors.New("HTTP header not set")
	}

	res, err = client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return nil, 500, errors.New("timeout")
		}
		return nil, 500, err
	}

	defer res.Body.Close()
	resp, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 500, err
	}
	statusCode = res.StatusCode

	switch statusCode {
	case 400, 500, 422, 401:
		errMsg := fmt.Sprintf("error from MyId: code %d msg %s", statusCode, resp)
		return resp, statusCode, errors.New(errMsg)

	case 200:
		return resp, statusCode, err

	default:
		errMsg := fmt.Sprintf("undefined code from MyId: code %d msg %s", statusCode, resp)
		return resp, statusCode, errors.New(errMsg)
	}
}
