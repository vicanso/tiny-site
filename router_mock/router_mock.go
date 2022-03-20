package routermock

import (
	"context"
	"encoding/json"
	"strings"

	"go.uber.org/atomic"

	"github.com/vicanso/tiny-site/email"
	"github.com/vicanso/tiny-site/log"
)

type (
	// RouterMock 路由配置信息
	RouterMock struct {
		Router     string `json:"router"`
		Route      string `json:"route"`
		Method     string `json:"method"`
		Status     int    `json:"status"`
		CotentType string `json:"cotentType"`
		Response   string `json:"response"`
		// DelaySeconds 延时，单位秒
		DelaySeconds int    `json:"delaySeconds"`
		URL          string `json:"url"`
	}
)

var currentRouterMocks = atomic.Value{}

// 更新router config配置
func Update(configs []string) {
	result := make(map[string]*RouterMock)
	for _, item := range configs {
		v := &RouterMock{}
		err := json.Unmarshal([]byte(item), v)
		if err != nil {
			log.Error(context.Background()).
				Err(err).
				Msg("router config is invalid")
			email.AlarmError(context.Background(), "router config is invalid:"+err.Error())
			continue
		}
		arr := strings.Split(v.Router, " ")
		if len(arr) == 2 {
			v.Method = arr[0]
			v.Route = arr[1]
		}
		// 如果未配置Route或者method的则忽略
		if v.Route == "" || v.Method == "" {
			continue
		}
		result[v.Method+v.Route] = v
	}

	currentRouterMocks.Store(result)
}

func getRouterMocks() map[string]*RouterMock {
	value := currentRouterMocks.Load()
	if value == nil {
		return nil
	}
	mocks, ok := value.(map[string]*RouterMock)
	if !ok {
		return nil
	}
	return mocks
}

// Get 获取路由配置
func Get(method, route string) *RouterMock {
	mocks := getRouterMocks()
	if mocks == nil {
		return nil
	}
	return mocks[method+route]
}

// List 获取路由mock配置
func List() map[string]RouterMock {
	result := make(map[string]RouterMock)
	for key, value := range getRouterMocks() {
		result[key] = *value
	}
	return result
}
