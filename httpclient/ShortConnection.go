package httpclient

import (
	"net/http"
	"io/ioutil"
	"time"
	"errors"
	"io"
	"strings"
)

/**
 * @desc 短连接请求
 * @author zhaojiangwei
 * @date 2018-12-25 17:59
 */

/**
 * @desc Get请求，请求结束后关闭连接
 * @author zhaojiangwei
 * @date 18:18 2018/12/25
 * @param timeout: 毫秒
 * @return
 **/
func ShortGet(uri string, timeout int, headers map[string]string, cookie []http.Cookie) (response ClientResponse) {
	uri = encodeUrl(uri)
	if uri == ""{
		return ClientResponse{Code: ClientNewReqFail, Error: errors.New("encode url error")}
	}

	response.Code = ClientCodeUnset
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		response.Code = ClientNewReqFail
		response.Error = err
		return
	}

	return shortExecute(req, timeout, headers, cookie)
}

func ShortPost(url string, body string, timeout int, headers map[string]string, cookie []http.Cookie) (ClientResponse) {
	url = encodeUrl(url)
	if url == ""{
		return ClientResponse{Code: ClientNewReqFail, Error: errors.New("encode url error")}
	}

	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return ClientResponse{Code: ClientNewReqFail, Error: err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return shortExecute(req, timeout, headers, cookie)
}

func shortExecute(req *http.Request, timeout int, headers map[string]string, cookie []http.Cookie) (response ClientResponse){
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	if cookie != nil && len(cookie) > 0 {
		for _, v := range cookie {
			req.AddCookie(&v)
		}
	}

	startTime := time.Now()
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	resp, resErr := client.Do(req)
	response.Cost = time.Now().Sub(startTime).Nanoseconds() / 1000 / 1000
	if resErr != nil {
		response.Code = ClientDoReqFail
		response.Error = resErr
		return
	}

	defer resp.Body.Close()

	response.Header = resp.Header
	if resp.StatusCode != 200 {
		response.Code = ClientHttpCodeFail
		response.StatusCode = resp.StatusCode
		response.Error = errors.New("status code exception")
		return
	}

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		response.Code = ClientReadBodyFail
		response.StatusCode = resp.StatusCode
		response.Error = bodyErr
		return
	}

	response.Code = ClientSucc
	response.StatusCode = resp.StatusCode
	response.Data = body
	return
}