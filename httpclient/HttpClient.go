package httpclient

import (
	"net/http"
	"net"
	"time"
	"io/ioutil"
	"io"
	"strings"
	"errors"
	"compress/gzip"
)

/**
 * @desc HttpClient 长连接使用
 * @author zhaojiangwei
 * @date 2018-12-25 14:44
 */

const(
	DefaultUserAgent = "GothicServer"
)

var (
	httpClient *http.Client
)

type HttpClient struct {
	client *http.Client
}

var DefaultPoolSetting = HttpPoolSetting{
	MaxIdleConns: 		10,
	IdleConnTimeout: 	60,
	HttpTimeout: 		5000,
	DialTimeout: 		2000,
}

var DefaultHttpSetting = HttpSetting{
	UserAgent:      DefaultUserAgent,
	Gzip:          	false,
	Headers:		nil,
	Cookies:		nil,
}

//连接池设置
type HttpPoolSetting struct {
	MaxIdleConns 		int
	IdleConnTimeout 	int		//单位s
	HttpTimeout 		int		//ms
	DialTimeout			int		//ms
}

//http请求设置
type HttpSetting struct {
	Gzip bool
	UserAgent string
	Headers map[string]string
	Cookies []http.Cookie
}

func NewDefaultHttpClient() *HttpClient {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(DefaultPoolSetting.DialTimeout) * time.Millisecond,
			}).DialContext,
			MaxIdleConns:        DefaultPoolSetting.MaxIdleConns,
			MaxIdleConnsPerHost: DefaultPoolSetting.MaxIdleConns,
			IdleConnTimeout:     time.Duration(DefaultPoolSetting.IdleConnTimeout) * time.Second,
		},

		Timeout: time.Duration(DefaultPoolSetting.HttpTimeout) * time.Millisecond,
	}

	return &HttpClient{
		client: client,
	}
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
func NewHttpClient(setting HttpPoolSetting) *HttpClient {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(setting.DialTimeout) * time.Millisecond,
			}).DialContext,
			MaxIdleConns:        setting.MaxIdleConns,
			MaxIdleConnsPerHost: setting.MaxIdleConns,
			IdleConnTimeout:     time.Duration(setting.IdleConnTimeout) * time.Second,
		},

		Timeout: time.Duration(setting.HttpTimeout) * time.Millisecond,
	}

	return &HttpClient{
		client: client,
	}
}

func (this *HttpClient) Get(uri string, setting *HttpSetting) (*ClientResponse) {
	uri = encodeUrl(uri)
	if uri == ""{
		return &ClientResponse{Code: ClientNewReqFail, Error: errors.New("encode url error")}
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return &ClientResponse{Code: ClientNewReqFail, Error: err}
	}

	return this.execute(req, setting)
}

func (this *HttpClient) Post(url string, body string, setting *HttpSetting) (*ClientResponse) {
	url = encodeUrl(url)
	if url == ""{
		return &ClientResponse{Code: ClientNewReqFail, Error: errors.New("encode url error")}
	}

	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return &ClientResponse{Code: ClientNewReqFail, Error: err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return this.execute(req, setting)
}

func (this *HttpClient) execute(req *http.Request, setting *HttpSetting) (response *ClientResponse) {
	response = &ClientResponse{}
	response.Code = ClientCodeUnset

	if setting == nil{
		setting = &DefaultHttpSetting
	}

	if setting.Headers != nil {
		for k, v := range setting.Headers {
			req.Header.Add(k, v)
		}
	}

	if setting.UserAgent != ""{
		req.Header.Set("User-Agent", setting.UserAgent)
	}else{
		req.Header.Set("User-Agent", DefaultUserAgent)
	}

	if setting.Cookies != nil && len(setting.Cookies) > 0 {
		for _, v := range setting.Cookies {
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

	response.Header = resp.Header
	if resp.StatusCode != 200 {
		response.Code = ClientHttpCodeFail
		response.StatusCode = resp.StatusCode
		response.Error = errors.New("status code exception")
		return
	}

	var body []byte
	var bodyErr error
	if setting.Gzip && resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			response.Code = ClientReadBodyFail
			response.StatusCode = resp.StatusCode
			response.Error = err
			return
		}

		body, bodyErr = ioutil.ReadAll(reader)
	}else{
		body, bodyErr = ioutil.ReadAll(resp.Body)
	}

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
