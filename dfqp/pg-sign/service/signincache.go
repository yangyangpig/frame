package service

import (
	"time"
	"fmt"
)

const (
	REDIS_MASTER_PREFIX = "signinrecord" // redis key 前缀
)
type signinCacheService struct {
	redisService
}

func (signincache *signinCacheService) getKey(mid int64) string {

	format := time.Now().Format("0102")
	return fmt.Sprintf("%s:%d:%d", REDIS_MASTER_PREFIX, format, mid)
}

