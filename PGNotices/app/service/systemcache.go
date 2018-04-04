package service

import (
	"time"
	"reflect"
)

const PREFIX = "syscache_"
type systemCacheService struct {}

func (c *systemCacheService) Get(key string) (interface{}, error) {
	ckey := c.getKey(key)
	mem, err := NewMemcache("notices")
	if err != nil {
		return "", err
	}
	res := mem.Get(ckey)
	switch res.(type) {
	case error:
		return nil, res.(error)
	case []byte:
		return string(res.([]byte)), nil
	case string:
		return reflect.ValueOf(res).String(), nil
	case float32, float64:
		return reflect.ValueOf(res).Float(), nil
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(res).Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(res).Uint(), nil
	default:
		return nil, nil
	}
}

func (c *systemCacheService) Set(key string, value interface{}, expire time.Duration) bool {
	ckey := c.getKey(key)
	mem, err := NewMemcache("notices")
	if err != nil {
		return false
	}
	errr := mem.Put(ckey, value, expire)
	if errr != nil {
		return false
	}
	return true
}

func (c *systemCacheService) getKey(key string) string {
	return PREFIX + key
}


