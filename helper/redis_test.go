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

package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisStats(t *testing.T) {
	assert := assert.New(t)
	m := RedisStats()
	assert.NotNil(m)
}

func TestRedisPing(t *testing.T) {
	assert := assert.New(t)
	err := RedisPing()
	assert.Nil(err)
}
