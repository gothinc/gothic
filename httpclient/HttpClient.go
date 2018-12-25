package httpclient

import (
	"net/http"
	"net"
	"time"
	"io/ioutil"
	"net/url"
	"io"
	"strings"
	"errors"
)

/**
 * @desc HttpClient 长连接使用
 * @author zhaojiangwei
 * @date 2018-12-25 14:44
 */

var (
	httpClient *http.Client
)

type HttpClient struct {
	client *http.Client
}

/**
 * @desc
 * @author zhaojiangwei
 * @date 15:10 2018/12/25
 * @param idleConnTimeout: 连接最长空闲时间
 * @param dialTimeout: 建立tcp连接超时时间
 * @param dialKeepalive: 长连接空闲时间
 * @return
 **/
func NewHttpClient(maxIdleConns, idleConnTimeout, httpTimeout, dialTimeout int) *HttpClient {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(dialTimeout) * time.Second,
			}).DialContext,
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConns,
			IdleConnTimeout:     time.Duration(idleConnTimeout) * time.Second,
		},

		Timeout: time.Duration(httpTimeout) * time.Millisecond,
	}

	return &HttpClient{
		client: client,
	}
}

func (this *HttpClient) Get(uri string, headers map[string]string, cookie []http.Cookie) (ClientResponse) {
	uri = this.encodeUrl(uri)
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return ClientResponse{Code: ClientNewReqFail, Error: err}
	}

	return this.execute(req, headers, cookie)
}

func (this *HttpClient) Post(url string, body string, headers map[string]string, cookie []http.Cookie) (ClientResponse) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return ClientResponse{Code: ClientNewReqFail, Error: err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return this.execute(req, headers, cookie)
}

func (this *HttpClient) execute(req *http.Request, headers map[string]string, cookie []http.Cookie) (response ClientResponse) {
	response.Code = ClientCodeUnset
	
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
	resp, resErr := this.client.Do(req)
	response.Cost = time.Now().Sub(startTime).Nanoseconds() / 1000 / 1000
	if resErr != nil {
		response.Code = ClientDoReqFail
		response.Error = resErr
		return
	}

	defer resp.Body.Close()
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

func (this *HttpClient) encodeUrl(uri string) string {
	urlObj, _ := url.Parse(uri)
	queryObj := urlObj.Query()
	urlObj.RawQuery = queryObj.Encode()
	return urlObj.String()
}
