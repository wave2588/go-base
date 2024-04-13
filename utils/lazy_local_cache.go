package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type lazyLocalCacheItem struct {
	value    interface{} // 缓存值
	expireAt int64       // 过期时间，秒级时间戳
}

type lazyLocalCacheImpl struct {
	itemMap map[string]*lazyLocalCacheItem
	rwMutex sync.RWMutex
}

type LazyLocalCacheLoadFunc = func(context.Context) (interface{}, error)

type lazyLocalCacheLoadTask struct {
	loadFunc LazyLocalCacheLoadFunc
	item     *lazyLocalCacheItem
	key      string
	ttl      int64
}

var (
	defaultLazyLocalCache = &lazyLocalCacheImpl{
		itemMap: make(map[string]*lazyLocalCacheItem, 100),
		rwMutex: sync.RWMutex{},
	}
	registeredLoadTaskList = make([]*lazyLocalCacheLoadTask, 0, 100)
)

func init() {
	go func() {
		for {
			now := CurrentSecond()
			for _, task := range registeredLoadTaskList {
				// 提前 2 秒准备刷新缓存
				if now+2 > task.item.expireAt {
					value, err := task.Load()
					if err != nil {
					} else {
						task.item.value = value
						task.item.expireAt = CurrentSecond() + task.ttl
					}
				}
			}
			time.Sleep(time.Second)
		}
	}()
}

func (t *lazyLocalCacheImpl) GetItem(key string) *lazyLocalCacheItem {
	t.rwMutex.RLock()
	result := t.itemMap[key]
	t.rwMutex.RUnlock()
	return result
}

func (t *lazyLocalCacheImpl) SetItemNX(key string, item *lazyLocalCacheItem) bool {
	defer t.rwMutex.Unlock()
	t.rwMutex.Lock()
	if t.itemMap[key] != nil {
		return false
	}
	t.itemMap[key] = item
	return true
}

func (t *lazyLocalCacheImpl) DeleteItem(key string) {
	t.rwMutex.Lock()
	delete(t.itemMap, key)
	t.rwMutex.Unlock()
}

func (t *lazyLocalCacheLoadTask) Load() (result interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			if pErr, ok := p.(error); ok {
				err = pErr
			} else {
				err = fmt.Errorf("%+v", p)
			}
		}
	}()
	result, err = t.loadFunc(context.Background())
	return
}

func GetLazyLocalCache(ctx context.Context, key string, loadFunc LazyLocalCacheLoadFunc, expireSeconds int64) (interface{}, error) {
	now := CurrentSecond()
	item := defaultLazyLocalCache.GetItem(key)
	if item == nil {
		value, err := loadFunc(ctx)
		if err != nil {
			return nil, err
		}
		item = &lazyLocalCacheItem{
			value:    value,
			expireAt: now + expireSeconds,
		}
		if defaultLazyLocalCache.SetItemNX(key, item) {
			registeredLoadTaskList = append(registeredLoadTaskList, &lazyLocalCacheLoadTask{
				loadFunc: loadFunc,
				item:     item,
				key:      key,
				ttl:      expireSeconds,
			})
		}
		return value, nil
	}
	return item.value, nil
}

func DeleteLazyLocalCache(key string) {
	defaultLazyLocalCache.DeleteItem(key)
}
