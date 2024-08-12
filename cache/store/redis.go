package store

import (
	"context"
	"reflect"
	"time"

	"github.com/wave2588/go-base/redis"
	"github.com/wave2588/go-base/utils"
)

type RedisStore struct {
	pool redis.Pool

	serializer Serializer
}

// NewRedisStore set default behavior as:
//
//		compress:   false
//	 escapeHTML: true
//	 serializer: JSON
func NewRedisStore(pool redis.Pool, options ...func(*RedisStore)) Store {
	r := &RedisStore{
		pool:       pool,
		serializer: NewJsonSerializer(false, true),
	}
	for _, opt := range options {
		opt(r)
	}
	return r
}

func WithSerializer(t Serializer) func(*RedisStore) {
	return func(r *RedisStore) {
		r.serializer = t
	}
}

func (r *RedisStore) Get(ctx context.Context, key string, dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return ErrBadDstType
	}

	value, err := r.pool.Slave().Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return r.serializer.Decode([]byte(value), dst)
}

func (r *RedisStore) MustGet(ctx context.Context, key string, dst interface{}) {
	utils.PanicIf(r.Get(ctx, key, dst))
}

// GetMulti get the provided keys from redis and store it in dst.
// Because golang has no generic type, so result must be provided in params.
// dst must be a map or pointer-to-map.
func (r *RedisStore) GetMulti(ctx context.Context, keys []string, dstMap interface{}) error {
	dstPtrV := reflect.ValueOf(dstMap)
	dstV := reflect.Indirect(dstPtrV)
	if dstV.Kind() != reflect.Map {
		return ErrBadDstMapType
	}

	// nil map
	if dstPtrV.Kind() != reflect.Ptr && dstV.IsNil() {
		return ErrBadDstMapValue
	}

	if dstPtrV.Kind() == reflect.Ptr && dstV.IsNil() {
		m := reflect.MakeMap(reflect.MapOf(dstV.Type().Key(), dstV.Type().Elem()))
		dstV.Set(m)
	}

	values, err := r.pool.Slave().MGet(ctx, keys...).Result()
	if err != nil {
		return err
	}

	for i, value := range values {
		if value == nil {
			continue
		}

		v := reflect.New(dstV.Type().Elem())
		if v.Kind() != reflect.Ptr {
			v = v.Addr()
		}

		if err := r.serializer.Decode([]byte(value.(string)), v.Interface()); err != nil {
			return err
		}

		dstV.SetMapIndex(reflect.ValueOf(keys[i]), v.Elem())
	}

	return nil
}

func (r *RedisStore) MustGetMulti(ctx context.Context, keys []string, dstMap interface{}) {
	utils.PanicIf(r.GetMulti(ctx, keys, dstMap))
}

func (r *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.pool.Slave().Exists(ctx, key).Result()

	return result != 0, err
}

func (r *RedisStore) MustExists(ctx context.Context, key string) bool {
	ret, err := r.Exists(ctx, key)
	utils.PanicIf(err)
	return ret
}

func (r *RedisStore) ExistsMulti(ctx context.Context, keys ...string) ([]bool, error) {
	if len(keys) == 0 {
		return []bool{}, nil
	}

	existsCmd := make([]*redis.IntCmd, len(keys))
	ret := make([]bool, len(keys))
	_, err := r.pool.Slave().Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for i, key := range keys {
			existsCmd[i] = pipeliner.Exists(ctx, key)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	for i, cmd := range existsCmd {
		ret[i] = cmd.Val() != 0
	}

	return ret, nil
}

func (r *RedisStore) MustExistsMulti(ctx context.Context, keys ...string) []bool {
	ret, err := r.ExistsMulti(ctx, keys...)
	utils.PanicIf(err)
	return ret
}

func (r *RedisStore) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	buf, err := r.serializer.Encode(value)
	if err != nil {
		return err
	}

	var v interface{} = buf
	_, err = r.pool.Master().Set(ctx, key, v, ttl).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) MustSet(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	utils.PanicIf(r.Set(ctx, key, value, ttl))
}

func (r *RedisStore) SetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration) error {
	srcV := reflect.Indirect(reflect.ValueOf(values))
	if srcV.Kind() != reflect.Slice {
		return ErrBadSrcType
	}
	if srcV.Len() != len(keys) {
		return ErrKeysLengthNotMatch
	}

	statusCmds := make([]*redis.StatusCmd, len(keys))
	_, err := r.pool.Master().Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for index, key := range keys {
			v := srcV.Index(index)
			if v.Kind() != reflect.Ptr {
				v = v.Addr()
			}

			if buf, err := r.serializer.Encode(v.Interface()); err != nil {
				return err
			} else {
				var v interface{} = buf
				statusCmds[index] = pipeliner.Set(ctx, key, v, ttl)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) MustSetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration) {
	utils.PanicIf(r.SetMulti(ctx, keys, values, ttl))
}

func (r *RedisStore) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	err := r.pool.Master().Del(ctx, keys...).Err()

	return err
}

func (r *RedisStore) MustDelete(ctx context.Context, keys ...string) {
	utils.PanicIf(r.Delete(ctx, keys...))
}
