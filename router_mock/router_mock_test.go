package routermock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterConfig(t *testing.T) {
	assert := assert.New(t)
	Update([]string{
		`{
			"route": "/",
			"method": "GET",
			"status": 400
		}`,
	})
	routeConfig := Get("GET", "/")
	assert.Equal(400, routeConfig.Status)
}
