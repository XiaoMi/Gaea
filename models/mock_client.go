package models

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockClient 是 Client 接口的 Mock 实现
type MockClient struct {
	mock.Mock
}

// Create 模拟 Create 方法
func (m *MockClient) Create(path string, data []byte) error {
	args := m.Called(path, data)
	return args.Error(0)
}

// Update 模拟 Update 方法
func (m *MockClient) Update(path string, data []byte) error {
	args := m.Called(path, data)
	return args.Error(0)
}

// UpdateWithTTL 模拟 UpdateWithTTL 方法
func (m *MockClient) UpdateWithTTL(path string, data []byte, ttl time.Duration) error {
	args := m.Called(path, data, ttl)
	return args.Error(0)
}

// Delete 模拟 Delete 方法
func (m *MockClient) Delete(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

// Read 模拟 Read 方法
func (m *MockClient) Read(path string) ([]byte, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}

// List 模拟 List 方法
func (m *MockClient) List(path string) ([]string, error) {
	args := m.Called(path)
	return args.Get(0).([]string), args.Error(1)
}

// ListWithValues 模拟 ListWithValues 方法
func (m *MockClient) ListWithValues(path string) (map[string]string, error) {
	args := m.Called(path)
	return args.Get(0).(map[string]string), args.Error(1)
}

// Close 模拟 Close 方法
func (m *MockClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// BasePrefix 模拟 BasePrefix 方法
func (m *MockClient) BasePrefix() string {
	args := m.Called()
	return args.String(0)
}
