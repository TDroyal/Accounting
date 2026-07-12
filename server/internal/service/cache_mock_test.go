// Package service_test 单测。用本地内存 mock 替代真实 Redis/DB，
// 覆盖统计聚合、缓存失效、分类禁用校验、JWT 往返等核心逻辑。
package service_test

import (
	"context"
	"sync"
	"time"
)

// memCache 内存实现的 cache.Cache，用于单测，无需真实 Redis。
type memCache struct {
	mu      sync.Mutex
	data    map[string]string
	sessions map[string]bool
}

func newMemCache() *memCache {
	return &memCache{data: map[string]string{}, sessions: map[string]bool{}}
}

func (m *memCache) SetSession(_ context.Context, userID, jti string, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[userID+":"+jti] = true
	return nil
}
func (m *memCache) DelSession(_ context.Context, userID, jti string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, userID+":"+jti)
	return nil
}
func (m *memCache) HasSession(_ context.Context, userID, jti string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sessions[userID+":"+jti], nil
}
func (m *memCache) GetStat(_ context.Context, key string) (string, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	return v, ok, nil
}
func (m *memCache) SetStat(_ context.Context, key, value string, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}
func (m *memCache) DelStat(_ context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}
func (m *memCache) Get(_ context.Context, key string) (string, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	return v, ok, nil
}
func (m *memCache) Set(_ context.Context, key, value string, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}
func (m *memCache) Del(_ context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}

// keyWasDeleted 用于断言缓存被失效（单测检查 DelStat 是否被调用）。
func (m *memCache) exists(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.data[key]
	return ok
}
