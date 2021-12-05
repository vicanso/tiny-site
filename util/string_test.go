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

package util

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(10, len(RandomString(10)))

	reg := regexp.MustCompile(`\d{10}`)
	assert.True(reg.MatchString(RandomDigit(10)))
}

func TestGenXID(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(20, len(GenXID()))
}

func TestSha256(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("ungWv48Bz+pBQUDeXa4iI7ADYaOWF3qctBD/YfIAFa0=", Sha256("abc"))
}

func TestContainsString(t *testing.T) {
	assert := assert.New(t)
	assert.True(ContainsString([]string{
		"a",
		"b",
	}, "a"))
	assert.False(ContainsString([]string{
		"a",
		"b",
	}, "c"))
}

func TestContainsAny(t *testing.T) {
	assert := assert.New(t)
	assert.True(ContainsAny([]string{
		"a",
		"b",
	}, []string{
		"a",
	}))
	assert.False(ContainsAny([]string{
		"a",
		"b",
	}, []string{
		"c",
	}))
}

func TestEncrypt(t *testing.T) {
	assert := assert.New(t)
	key := []byte("01234567890123456789012345678901")
	data := []byte("abcd")

	result, err := Encrypt(key, data)
	assert.Nil(err)
	result, err = Decrypt(key, result)
	assert.Nil(err)
	assert.Equal(data, result)
}

func TestGetFirstLetter(t *testing.T) {
	assert := assert.New(t)
	value := GetFirstLetter("测试")
	assert.Equal("C", value)
}
