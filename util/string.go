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
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/mozillazg/go-pinyin"
	"github.com/rs/xid"
)

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const digitBytes = "0123456789"

// randomString 生成随机的字符串
func randomString(baseLetters string, n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(baseLetters) {
			b[i] = baseLetters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandomString 创建指定长度的字符串
func RandomString(n int) string {
	return randomString(letterBytes, n)
}

// RandomDigit 创建指定长度的数字字符串
func RandomDigit(n int) string {
	return randomString(digitBytes, n)
}

// GenXID 生成xid
func GenXID() string {
	return strings.ToUpper(xid.New().String())
}

// Sha256 对字符串做sha256后返回base64字符串
func Sha256(str string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashBytes)
}

// ContainsString 判断字符串数组是否包含该字符串
func ContainsString(arr []string, str string) (found bool) {
	for _, v := range arr {
		if found {
			break
		}
		if v == str {
			found = true
		}
	}
	return
}

// ContainsAny 判断该字符串数据是否包含其中任意一个字符串
func ContainsAny(targets []string, checkArr []string) bool {
	valid := false
	for _, item := range targets {
		if ContainsString(checkArr, item) {
			valid = true
			break
		}
	}
	return valid
}

// Encrypt 数据加密
// https://stackoverflow.com/questions/18817336/golang-encrypting-a-string-with-aes-and-base64
func Encrypt(key, text []byte) ([]byte, error) {
	// 需要注意 key的长度必须为32字节
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// Decrypt 数据解密
// https://stackoverflow.com/questions/18817336/golang-encrypting-a-string-with-aes-and-base64
func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetFirstLetter 获取字符串首字母
func GetFirstLetter(str string) string {
	arr := pinyin.LazyPinyin(str, pinyin.NewArgs())
	if len(arr) == 0 {
		return strings.ToUpper(str[0:1])
	}
	return strings.ToUpper(arr[0][0:1])
}

// CutRune 按rune截断字符串
func CutRune(str string, max int) string {
	result := []rune(str)
	if len(result) < max {
		return str
	}
	return string(result[:max]) + "..."
}
