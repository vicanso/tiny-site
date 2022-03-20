// Copyright 2020 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package location

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/go-axios"
)

func TestGetLocationByIP(t *testing.T) {
	assert := assert.New(t)

	done := ins.Mock(&axios.Response{
		Status: 200,
		Data: []byte(`{
			"ip": "47.242.95.102",
			"country": "美国",
			"province": "加利福尼亚",
			"city": "圣克拉拉",
			"isp": "阿里巴巴"
		  }`),
	})
	defer done()
	lo, err := GetByIP(context.Background(), "127.0.0.1")
	assert.Nil(err)
	assert.Equal("圣克拉拉", lo.City)
}
