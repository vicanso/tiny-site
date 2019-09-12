// Copyright 2019 tree xie
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
	"sync/atomic"
	"unsafe"

	"github.com/vicanso/ips"
)

var (
	ipBlocker = &IPBlocker{
		IPS: ips.New(),
	}
)

type (
	// IPBlocker ip blocker
	IPBlocker struct {
		IPS *ips.IPS
	}
)

// ResetIPBlocker reset the ip blocker
func ResetIPBlocker(ipList []string) {
	iPS := ips.New()
	for _, value := range ipList {
		iPS.Add(value)
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&ipBlocker.IPS)), unsafe.Pointer(iPS))
}

// IsBlockIP check the ip is blocked
func IsBlockIP(ip string) bool {
	iPS := (*ips.IPS)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&ipBlocker.IPS))))
	blocked := iPS.Exists(ip)
	return blocked
}
