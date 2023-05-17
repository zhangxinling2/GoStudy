package cache

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMaxMemoryCache_Set(t *testing.T) {
	testCases := []struct {
		name     string
		cache    func() *MaxMemoryCache
		key      string
		val      []byte
		wantErr  error
		wantKey  string
		wantUsed int64
	}{
		{
			name: "set success with not exist",
			cache: func() *MaxMemoryCache {
				res := NewMaxMemoryCache(100, &mockCache{
					data: map[string][]byte{},
				})
				return res
			},
			key:      "key1",
			val:      []byte("hello"),
			wantUsed: 5,
			wantKey:  "key1",
		},
		{
			name: "set success with already exist",
			cache: func() *MaxMemoryCache {
				res := NewMaxMemoryCache(100, &mockCache{
					data: map[string][]byte{"key1": []byte("hello")},
				})
				return res
			},
			key:      "key1",
			val:      []byte("hello,key2"),
			wantUsed: 10,
			wantKey:  "key1",
		},
		{
			name: "LRU DELETE TEST",
			cache: func() *MaxMemoryCache {
				res := NewMaxMemoryCache(14, &mockCache{
					data: map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
				})
				res.Set(context.Background(), "key2", []byte("value2"), time.Minute)
				res.Set(context.Background(), "key1", []byte("value1"), time.Minute)
				return res
			},
			key:      "key3",
			val:      []byte("hello"),
			wantUsed: 11,
			wantKey:  "key3",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := tc.cache()
			err := cache.Set(context.Background(), tc.key, tc.val, time.Minute)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantUsed, cache.used)
			assert.Equal(t, tc.wantKey, cache.cache.head.Next.key)
		})
	}
}
func TestMaxMemoryCache_Get(t *testing.T) {
	testCases := []struct {
		name     string
		cache    func() *MaxMemoryCache
		key      string
		val      any
		wantErr  error
		wantKey  string
		wantUsed int64
	}{
		{
			name: "get success",
			cache: func() *MaxMemoryCache {
				res := NewMaxMemoryCache(100, &mockCache{
					data: map[string][]byte{
						"key1": []byte("value1"),
						"key2": []byte("value2"),
						"key3": []byte("value3"),
					},
				})
				res.Set(context.Background(), "key3", []byte("value3"), time.Minute)
				res.Set(context.Background(), "key2", []byte("value2"), time.Minute)
				res.Set(context.Background(), "key1", []byte("value1"), time.Minute)
				return res
			},
			key:      "key3",
			wantUsed: 18,
			wantKey:  "key3",
		},
		{
			name: "no key",
			cache: func() *MaxMemoryCache {
				res := NewMaxMemoryCache(100, &mockCache{
					data: map[string][]byte{},
				})
				return res
			},
			key:     "key1",
			wantErr: errors.New("key不存在"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := tc.cache()
			_, err := cache.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantUsed, cache.used)
			assert.Equal(t, tc.wantKey, cache.cache.head.Next.key)
		})
	}
}

type mockCache struct {
	fn   func(key string, val []byte)
	data map[string][]byte
}

func (m *mockCache) Set(ctx context.Context, key string, val []byte, expiration time.Duration) error {
	m.data[key] = val
	return nil
}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, ok := m.data[key]
	if ok {
		return val, nil
	}
	return nil, errors.New("key不存在")
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	val, ok := m.data[key]
	if ok {
		m.fn(key, val)
	}
	return nil
}

func (m *mockCache) LoadAndDelete(ctx context.Context, key string) ([]byte, error) {
	val, ok := m.data[key]
	if ok {
		m.fn(key, val)
		return val, nil
	}
	return nil, errors.New("key不存在")
}

func (m *mockCache) OnEvicted(fn func(key string, any []byte)) {
	m.fn = fn
}
