package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseUrl       string
	Authorization string
	Debug         bool
}

func (c *Client) Get(path string, v any) *Error {
	return c.doRequest(path, http.MethodGet, nil, v)
}

func (c *Client) Post(path string, req any, v any) *Error {
	return c.doRequest(path, http.MethodPost, req, v)
}

func (c *Client) Patch(path string, req any, v any) *Error {
	return c.doRequest(path, http.MethodPost, req, v)
}

func (c *Client) Delete(path string, v any) *Error {
	return c.doRequest(path, http.MethodDelete, nil, v)
}

func (c *Client) doRequest(path string, method string, req any, v any) *Error {
	var reqBody = new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(req)

	var hc http.Client
	request, err := http.NewRequest(method, c.BaseUrl+path, reqBody)
	if err != nil {
		return &Error{StatusCode: 0, ErrCode: ErrCodeNetwork, Message: err.Error()}
	}
	request.Header.Set("Content-Type", "application/json")
	if c.Authorization != "" {
		request.Header.Set("Authorization", c.Authorization)
	}

	res, err := hc.Do(request)
	if err != nil {
		return &Error{StatusCode: 0, ErrCode: ErrCodeNetwork, Message: err.Error()}
	}
	var buf = new(bytes.Buffer)
	io.Copy(buf, res.Body)
	defer func() {
		res.Body.Close()
		if c.Debug {
			fmt.Printf("%s: %s, req: %s, resp: %s\n", method, path, reqBody.String(), buf.String())
		}
	}()

	if res.StatusCode > http.StatusCreated {
		return newClientErrorFromResponse(res.StatusCode, buf.Bytes())
	}
	err = json.Unmarshal(buf.Bytes(), v)
	if err != nil {
		return &Error{StatusCode: 0, ErrCode: ErrCodeBadResponseBody, Message: err.Error()}
	}
	return nil
}

func newClientErrorFromResponse(statusCode int, body []byte) *Error {
	var ce = &Error{
		StatusCode: statusCode,
		ErrCode:    ErrCodeOk,
		Message:    "",
	}

	if statusCode > http.StatusCreated {
		err2 := json.Unmarshal(body, ce)
		if err2 != nil {
			ce.ErrCode = ErrCodeUnknown
			ce.Message = string(body)
			return ce
		}

		if ce.Message == "" {
			ce.Message = string(body)
		}
	}

	return ce
}
