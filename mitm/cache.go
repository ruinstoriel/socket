// Copyright 2018 ouqiang authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package cert HTTPS证书
package main

import (
	"crypto/tls"
	"sync"
)

// Cache 证书缓存接口
type Cache interface {
	Set(host string, c *tls.Certificate)
	Get(host string) *tls.Certificate
}

// 实现证书缓存接口
type CacheImp struct {
	m sync.Map
}

func (c *CacheImp) Set(host string, cert *tls.Certificate) {
	c.m.Store(host, cert)
}
func (c *CacheImp) Get(host string) *tls.Certificate {
	v, ok := c.m.Load(host)
	if !ok {
		return nil
	}

	return v.(*tls.Certificate)
}
