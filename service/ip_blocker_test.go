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

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPBlocker(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(ResetIPBlocker([]string{
		"1.1.1.1",
	}))
	assert.True(IsBlockIP("1.1.1.1"))
	assert.False(IsBlockIP("1.1.1.2"))
}
