package entity

import "fmt"

const userBagCacheKey = "bag|%d"
const userNewGoodsCacheKey = "bag|newItems|%d"



// UserBagCacheKey 用户背包cache key
func UserBagCacheKey(mid int64) string {
	return fmt.Sprintf(userBagCacheKey, mid)
}

// UserNewGoodsCacheKey 背包中新物品cache key
func UserNewGoodsCacheKey(mid int64) string {
	return fmt.Sprintf(userNewGoodsCacheKey, mid)
}
