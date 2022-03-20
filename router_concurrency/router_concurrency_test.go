package routerconcurrency

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestRouterLimiter(t *testing.T) {
	assert := assert.New(t)

	InitLimiter([]elton.RouterInfo{
		{
			Method: "GET",
			Route:  "/",
		},
		{
			Method: "GET",
			Route:  "/users/me",
		},
	})
	rc := GetLimiter()
	key := "GET /"
	count, max := rc.IncConcurrency(key)
	assert.Equal(uint32(1), count)
	assert.Equal(uint32(0), max)

	assert.Equal(uint32(1), rc.GetConcurrency(key))

	rc.DecConcurrency(key)
	assert.Equal(uint32(0), rc.GetConcurrency(key))

	// 重置路由并发配置
	Update([]string{
		`{
			"router": "GET /",
			"max": 10,
			"rate": 100,
			"interval": "1s"
		}`,
	})
	count, max = rc.IncConcurrency(key)
	assert.Equal(uint32(1), count)
	assert.Equal(uint32(10), max)
}
