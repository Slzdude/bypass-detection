package models

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"github.com/Skactor/bypass-detection/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func ParseUrl(url *url.URL) *UrlType {
	return &UrlType{
		Scheme:   url.Scheme,
		Domain:   url.Hostname(),
		Host:     url.Host,
		Port:     url.Port(),
		Path:     url.EscapedPath(),
		Query:    url.RawQuery,
		Fragment: url.Fragment,
	}
}

func ParseRequest(oReq *http.Request) (*Request, error) {
	req := &Request{}
	req.Method = oReq.Method
	req.Url = ParseUrl(oReq.URL)
	header := make(map[string]string)
	for k := range oReq.Header {
		header[k] = oReq.Header.Get(k)
	}
	req.Headers = header
	req.ContentType = oReq.Header.Get("Content-Type")
	if oReq.Body == nil || oReq.Body == http.NoBody {
	} else {
		data, err := ioutil.ReadAll(oReq.Body)
		if err != nil {
			return nil, err
		}
		req.Body = data
		oReq.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}
	return req, nil
}

func ParseResponse(oResp *http.Response) (*Response, error) {
	var resp Response
	header := make(map[string]string)
	resp.Status = int32(oResp.StatusCode)
	resp.Url = ParseUrl(oResp.Request.URL)
	for k := range oResp.Header {
		header[k] = oResp.Header.Get(k)
	}
	resp.Headers = header
	resp.ContentType = oResp.Header.Get("Content-Type")
	body, err := getRespBody(oResp)
	if err != nil {
		return nil, err
	}
	resp.Body = body
	return &resp, nil
}

func getRespBody(oResp *http.Response) ([]byte, error) {
	var body []byte
	if oResp.Header.Get("Content-Encoding") == "gzip" {
		gr, _ := gzip.NewReader(oResp.Body)
		defer gr.Close()
		for {
			buf := make([]byte, 1024)
			n, err := gr.Read(buf)
			if err != nil && err != io.EOF {
				//utils.Logger.Error(err)
				return nil, err
			}
			if n == 0 {
				break
			}
			body = append(body, buf...)
		}
	} else {
		raw, err := ioutil.ReadAll(oResp.Body)
		if err != nil {
			//utils.Logger.Error(err)
			return nil, err
		}
		//defer oResp.Body.Close()
		body = raw
	}
	return body, nil
}

type Connection struct {
	Request  *Request
	Response *Response
}

func ReadHTTPFromBytes(request []byte, response []byte) *Connection {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(request)))
	if err != nil {
		logger.Logger.Errorf("Failed to load request bytes: %s", err.Error())
		return nil
	}
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(response)), req)
	if err != nil {
		return nil
	}
	//save response body
	b := new(bytes.Buffer)
	io.Copy(b, resp.Body)
	resp.Body.Close()
	resp.Body = ioutil.NopCloser(b)

	parsedRequest, err := ParseRequest(req)
	if err != nil {
		logger.Logger.Errorf("Failed to convert go object to protobuf object")
		return nil
	}
	parsedResponse, err := ParseResponse(resp)
	if err != nil {
		logger.Logger.Errorf("Failed to convert go object to protobuf object")
		return nil
	}
	return &Connection{
		Request:  parsedRequest,
		Response: parsedResponse,
	}
}
