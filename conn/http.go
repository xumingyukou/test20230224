package conn

//http request 工具函数
import (
	"bytes"
	"clients/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

var (
	fastHttpClient = &fasthttp.Client{
		Name:                "goex-http-logger",
		MaxConnsPerHost:     16,
		MaxIdleConnDuration: 20 * time.Second,
		ReadTimeout:         10 * time.Second,
		WriteTimeout:        10 * time.Second,
	}
	socksDialer fasthttp.DialFunc
)

type HttpError struct {
	Code    int32 // http状态码
	Unknown bool  // 请求执行状态是否未知
}

func (e HttpError) Error() string {
	return fmt.Sprintf("error: code=%d  unknown=%t", e.Code, e.Unknown)
}

func NewHttpRequestWithFasthttp(client *http.Client, reqMethod, reqUrl, postData string, headers map[string]string) ([]byte, error) {
	logger.Logger.Debug("use fasthttp client")
	transport := client.Transport

	if transport != nil {
		if proxy, err := transport.(*http.Transport).Proxy(nil); err == nil && proxy != nil {
			proxyUrl := proxy.String()
			logger.Logger.Debug("proxy url: ", proxyUrl)
			if proxy.Scheme != "socks5" {
				logger.Logger.Error("fasthttp only support the socks5 proxy")
			} else if socksDialer == nil {
				socksDialer = fasthttpproxy.FasthttpSocksDialer(strings.TrimPrefix(proxyUrl, proxy.Scheme+"://"))
				fastHttpClient.Dial = socksDialer
			}
		}
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.SetMethod(reqMethod)
	req.SetRequestURI(reqUrl)
	req.SetBodyString(postData)

	err := fastHttpClient.Do(req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), nil
}

func NewHttpRequest(client *http.Client, reqType string, reqUrl string, postData string, requstHeaders map[string]string) ([]byte, error) {
	logger.Logger.Debugf("[%s] request url: %s", reqType, reqUrl)
	lib := os.Getenv("HTTP_LIB")
	if lib == "fasthttp" {
		return NewHttpRequestWithFasthttp(client, reqType, reqUrl, postData, requstHeaders)
	}

	req, _ := http.NewRequest(reqType, reqUrl, strings.NewReader(postData))
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	}

	for k, v := range requstHeaders {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, httpError(err, resp)
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData))
	}

	return bodyData, nil
}

func NewHttpRequestWithHeader(client *http.Client, reqType string, reqUrl string, postData string, requstHeaders map[string]string) ([]byte, http.Header, error) {
	logger.Logger.Debugf("[%s] request url: %s", reqType, reqUrl)
	lib := os.Getenv("HTTP_LIB")
	if lib == "fasthttp" {
		res, err := NewHttpRequestWithFasthttp(client, reqType, reqUrl, postData, requstHeaders)
		return res, nil, err
	}

	req, _ := http.NewRequest(reqType, reqUrl, strings.NewReader(postData))
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	}
	for k, v := range requstHeaders {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, httpError(err, resp)
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData))
	}

	return bodyData, resp.Header, nil
}

func HttpGet(client *http.Client, reqUrl string) (map[string]interface{}, error) {
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", nil)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Println(string(respData))
		return nil, err
	}
	return bodyDataMap, nil
}

func Request(client *http.Client, url, method string, header *http.Header, params url.Values) ([]byte, error) {
	var (
		req *http.Request
		rsp []byte
		err error
	)
	if method == "POST" {
		var payload *strings.Reader
		if strings.Contains(header.Get("Content-Type"), "json") {
			params1 := make(map[string]string)
			for key, value := range params {
				params1[key] = value[0]
			}
			dataByte, err := json.Marshal(params1)
			if err != nil {
				return nil, err
			}
			payload = strings.NewReader(string(dataByte))
		} else {
			payload = strings.NewReader(params.Encode())
		}
		req, _ = http.NewRequest("POST", url, payload)
	} else {
		if params == nil {
			req, _ = http.NewRequest(method, url, nil)
		} else {
			url += "?" + params.Encode()
			req, _ = http.NewRequest(method, url, nil)
		}
	}
	req.Header = *header
	response, err := client.Do(req)
	if err != nil {
		return nil, httpError(err, response)
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			return
		}
	}()
	rsp, err = ioutil.ReadAll(response.Body)
	//logger.Logger.Debug(url, method, params, string(resBytes))
	if err != nil {
		err = httpError(err, response)
		return nil, err
	}
	if response.StatusCode > 400 {
		fmt.Println("返回error", response.StatusCode)

		err = httpError(err, response)
		fmt.Println("返回err:", err)
		return nil, err
	}
	return rsp, err
}

func HuoBiRequest(client *http.Client, url, method string, header *http.Header, params interface{}) ([]byte, error) {
	var (
		req *http.Request
		rsp []byte
		err error
	)

	if method == "GET" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		post_params, _ := json.Marshal(&params)
		payload := bytes.NewReader(post_params)
		req, err = http.NewRequest(method, url, payload)
	}

	req.Header = *header
	response, err := client.Do(req)
	if err != nil {
		return nil, httpError(err, response)
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			return
		}
	}()
	rsp, err = ioutil.ReadAll(response.Body)
	a := string(rsp)
	fmt.Println(a)
	if err != nil {
		err = httpError(err, response)
		return nil, err
	}
	if response.StatusCode > 400 {
		fmt.Println("返回error", response.StatusCode)

		err = httpError(err, response)
		fmt.Println("返回err:", err)
		return nil, err
	}

	return rsp, err
}

func HuoBiHttpGet(url string) ([]byte, error) {
	// logger := perflogger.GetInstance()
	// logger.Start()

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)

	// logger.StopAndLog("GET", url)

	return result, err
}

func DoFTXRequest(client *http.Client, request *http.Request) (rsp []byte, err error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, httpError(err, response)
	}
	defer func() {
		err2 := response.Body.Close()
		if err2 != nil {
			return
		}
	}()
	rsp, err = ioutil.ReadAll(response.Body)
	//logger.Logger.Debug(url, method, params, string(resBytes))
	if err != nil {
		err = httpError(err, response)
		return
	}
	return
}

func BatchRequest(client *http.Client, url string, header *http.Header, params []map[string]string) (rsp []byte, err error) {
	var req *http.Request
	var payload *strings.Reader
	dataByte, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	payload = strings.NewReader(string(dataByte))
	req, _ = http.NewRequest("POST", url, payload)
	req.Header = *header
	response, err := client.Do(req)
	if err != nil {
		return nil, httpError(err, response)
	}
	defer func() {
		err2 := response.Body.Close()
		if err2 != nil {
			return
		}
	}()
	rsp, err = ioutil.ReadAll(response.Body)
	//logger.Logger.Debug(url, method, params, string(resBytes))
	if err != nil {
		err = httpError(err, response)
		return
	}
	return
}

func RequestWithHeader(client *http.Client, url, method string, header *http.Header, params url.Values) (*http.Header, []byte, error) {
	var req *http.Request
	if method == "POST" {
		var payload *strings.Reader
		if strings.Contains(header.Get("Content-Type"), "json") {
			params1 := make(map[string]string)
			for key, value := range params {
				params1[key] = value[0]
			}
			dataByte, err := json.Marshal(params1)
			if err != nil {
				return nil, nil, err
			}
			payload = strings.NewReader(string(dataByte))
		} else {
			payload = strings.NewReader(params.Encode())
		}
		//j := getOkexPostPara(params)
		//req, _ = http.NewRequest("POST", url, bytes.NewBuffer(j))
		req, _ = http.NewRequest("POST", url, payload)
	} else {
		url += "?" + params.Encode()
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header = *header
	response, err := client.Do(req)
	if err != nil {
		return nil, nil, httpError(err, response)
	}
	defer func() {
		err2 := response.Body.Close()
		if err2 != nil {
			return
		}
	}()
	rsp, err := ioutil.ReadAll(response.Body)
	//logger.Logger.Debug(url, method, params, string(resBytes))
	if err != nil {
		err = httpError(err, response)
		return nil, nil, err
	}
	return &response.Header, rsp, nil
}

func HttpGet2(client *http.Client, reqUrl string, headers map[string]string) (map[string]interface{}, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Println("respData", string(respData))
		return nil, err
	}
	return bodyDataMap, nil
}

func HttpGet3(client *http.Client, reqUrl string, headers map[string]string) ([]interface{}, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}
	//println(string(respData))
	var bodyDataMap []interface{}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Println("respData", string(respData))
		return nil, err
	}
	return bodyDataMap, nil
}

func HttpGet4(client *http.Client, reqUrl string, headers map[string]string, result interface{}) error {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respData, result)
	if err != nil {
		log.Printf("HttpGet4 - json.Unmarshal failed : %v, resp %s", err, string(respData))
		return err
	}

	return nil
}
func HttpGet5(client *http.Client, reqUrl string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func HttpPostForm(client *http.Client, reqUrl string, postData url.Values) ([]byte, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded"}
	return NewHttpRequest(client, "POST", reqUrl, postData.Encode(), headers)
}

func HttpPostForm2(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "POST", reqUrl, postData.Encode(), headers)
}

func HttpPostForm3(client *http.Client, reqUrl string, postData string, headers map[string]string) ([]byte, error) {
	return NewHttpRequest(client, "POST", reqUrl, postData, headers)
}

func HttpPostForm4(client *http.Client, reqUrl string, postData map[string]string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/json"
	data, _ := json.Marshal(postData)
	return NewHttpRequest(client, "POST", reqUrl, string(data), headers)
}

func HttpDeleteForm(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "DELETE", reqUrl, postData.Encode(), headers)
}

func HttpPut(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewHttpRequest(client, "PUT", reqUrl, postData.Encode(), headers)
}

func httpError(err error, response *http.Response) *HttpError {
	statusCode := 0
	if response != nil {
		statusCode = response.StatusCode
	}

	unknown := false
	if err != nil || statusCode > 500 {
		unknown = true
	}

	if !unknown {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			unknown = true
		}
	}

	if statusCode == 0 {
		statusCode = 502
	}

	return &HttpError{Code: int32(statusCode), Unknown: unknown}
}

// ok post格式化
func getOkexPostPara(params url.Values) []byte {
	j, err1 := json.Marshal(params)
	if err1 != nil {
		fmt.Println(err1)
		panic(err1)
	}
	a := strings.ReplaceAll(strings.ReplaceAll(string(j), "[", ""), "]", "")
	j = []byte(a)
	return j
}

type SlidingWindow struct {
	// 每 period ms 最多 limit 个请求
	limit  int64 //
	period int64 //

	reqTimes []int64
	lock     sync.Mutex
}

func (s *SlidingWindow) Allow() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now().UnixMilli()
	if int64(len(s.reqTimes)) < s.limit {
		s.reqTimes = append(s.reqTimes, now)
		return true
	}
	// len(s.reqTimes) == s.limit
	// 窗口内最早的请求 s.reqTimes[0]
	if now-s.reqTimes[0] <= s.period {
		return false
	}
	s.reqTimes = append(s.reqTimes[1:], now)
	return true
}

func (s *SlidingWindow) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.reqTimes = make([]int64, 0, s.limit)
}

func NewSlidingWindow(limit int64, period int64) *SlidingWindow {
	s := &SlidingWindow{
		limit:    limit,
		period:   period,
		reqTimes: make([]int64, 0, limit),
	}
	return s
}
