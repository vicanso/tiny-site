package util

import (
	"io/ioutil"
	"net/http"

	"github.com/mozillazg/request"
)

var requestHooks []request.Hook

type (
	hookConverError struct {
	}
)

func init() {
	AddHook(&hookConverError{})
}

func createRequestError(msg string, req *http.Request, resp *http.Response) (err *HTTPError) {
	status := http.StatusInternalServerError
	if resp != nil {
		status = resp.StatusCode
	}
	err = &HTTPError{
		StatusCode: status,
		Message:    msg,
		Category:   ErrCategoryRequset,
	}
	extra := make(map[string]interface{})
	url := req.URL
	extra["uri"] = url.RequestURI()
	extra["host"] = url.Host
	extra["method"] = req.Method
	err.Extra = extra
	return
}

func (h *hookConverError) BeforeRequest(req *http.Request) (resp *http.Response, err error) {
	return
}
func (h *hookConverError) AfterRequest(req *http.Request, resp *http.Response, err error) (newResp *http.Response, newErr error) {
	if err != nil {
		newErr = createRequestError(err.Error(), req, nil)
		return
	}
	statusCode := resp.StatusCode
	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		return
	}

	// 对于<200或者>=400的出错转换
	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		newErr = createRequestError(err.Error(), req, nil)
		return
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(buf, &m)
	if err != nil {
		newErr = createRequestError(err.Error(), req, nil)
		return
	}

	// 可对于不同服务的响应获取出错信息
	message, ok := m["message"].(string)
	if !ok {
		message = "未知异常"
	}
	he := createRequestError(message, req, nil)
	he.Extra["body"] = string(buf)
	newErr = he
	return
}

// AddHook 增加请求hook处理
func AddHook(hook request.Hook) {
	requestHooks = append(requestHooks, hook)
}

// HTTPGet http get request
func HTTPGet(url string, params map[string]string) (data []byte, err error) {
	c := new(http.Client)
	req := request.NewRequest(c)
	req.Hooks = requestHooks
	if params != nil {
		req.Params = params
	}
	resp, err := req.Get(url)
	if err != nil {
		return
	}
	data, err = resp.Content()
	defer resp.Body.Close()
	return
}

// HTTPPost http post request
func HTTPPost(url string, body map[string]interface{}, params map[string]string) (data []byte, err error) {
	c := new(http.Client)
	req := request.NewRequest(c)
	req.Hooks = requestHooks
	if body != nil {
		req.Json = body
	}
	if params != nil {
		req.Params = params
	}
	resp, err := req.Post(url)
	if err != nil {
		return
	}
	data, err = resp.Content()
	defer resp.Body.Close()
	return
}
