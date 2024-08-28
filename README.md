# go-base


## Cache 使用

### 初始化
```go
import (
    "github.com/go-redis/redis/v8"
    "github.com/wave2588/go-base/cache"
    redis2 "github.com/wave2588/go-base/redis"
)
    
var DefaultCache *cache.MemCache

func init() {
	masterAddr := &url.URL{
        Scheme: "",
        Host:   "",
	}
	masterAddr.User = url.UserPassword("", "password")
	defaultCachePool, err := redis2.GetPool("dump", masterAddr, nil)
	utils.PanicIf(err)

	// store.NewLocalStore()
	DefaultCache = cache.NewMemCache(store.NewRedisStore(defaultCachePool), 72*time.Hour)
}
```
### 使用
```go

type Account struct {
    ID int64 `boil:"id" json:"id" toml:"id" yaml:"id"`
}

const GetAccountCacheKey = "GetAccount"

func BatchGetCache(ctx context.Context, ids []int64) (map[int64]*Account, error) {
	resultMap := make(map[int64]*Account)
	ttl := 12 * time.Hour
	err := DefaultCache.GetMulti(ctx, ids,
		func(i interface{}) string {
			return store.NewKey(GetAccountCacheKey, i.(int64))
		},
		func(i interface{}) (interface{}, error) {
			ids := i.([]int64)
			return d.BatchGetByMySQL(ctx, ids)
		}, &resultMap, &ttl)
	return resultMap, err
}

```
