// Copyright 2021 tree xie
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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlmostEqual(t *testing.T) {
	assert := assert.New(t)
	assert.True(AlmostEqual(0.1+0.2+9e-7, 0.3))
	assert.False(AlmostEqual(0.2+1e-5, 0.2))
}

func TestToFixed(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("1.12", ToFixed(1.1221132, 2))
	assert.Equal("1.13", ToFixed(1.129, 2))
	assert.Equal("1.10", ToFixed(1.1, 2))
	assert.Equal("1", ToFixed(1.1, 0))
}
