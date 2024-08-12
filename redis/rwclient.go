package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// https://github.com/twitter/twemproxy/blob/master/notes/redis.md#redis-command-support

var _ Client = new(rwClientImpl)

func NewRWClient(pool Pool) Client {
	return &rwClientImpl{pool: pool}
}

type rwClientImpl struct {
	pool Pool
}

func (c *rwClientImpl) Pipeline() Pipeliner {
	return c.pool.Master().Pipeline()
}

func (c *rwClientImpl) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.pool.Master().Pipelined(ctx, fn)
}

func (c *rwClientImpl) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.pool.Master().TxPipelined(ctx, fn)
}

func (c *rwClientImpl) TxPipeline() Pipeliner {
	return c.pool.Master().TxPipeline()
}

func (c *rwClientImpl) Command(ctx context.Context) *CommandsInfoCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClientGetName(ctx context.Context) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Echo(ctx context.Context, message interface{}) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Ping(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Quit(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Del(ctx context.Context, keys ...string) *IntCmd {
	return c.pool.Master().Del(ctx, keys...)
}

func (c *rwClientImpl) Unlink(ctx context.Context, keys ...string) *IntCmd {
	return c.pool.Master().Unlink(ctx, keys...)
}

func (c *rwClientImpl) Dump(ctx context.Context, key string) *StringCmd {
	return c.pool.Slave().Dump(ctx, key)
}

func (c *rwClientImpl) Exists(ctx context.Context, keys ...string) *IntCmd {
	return c.pool.Slave().Exists(ctx, keys...)
}

func (c *rwClientImpl) Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd {
	return c.pool.Master().Expire(ctx, key, expiration)
}

func (c *rwClientImpl) ExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd {
	return c.pool.Master().ExpireAt(ctx, key, tm)
}

func (c *rwClientImpl) Keys(ctx context.Context, pattern string) *StringSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Move(ctx context.Context, key string, db int) *BoolCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ObjectRefCount(ctx context.Context, key string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ObjectEncoding(ctx context.Context, key string) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ObjectIdleTime(ctx context.Context, key string) *DurationCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Persist(ctx context.Context, key string) *BoolCmd {
	return c.pool.Master().Persist(ctx, key)
}

func (c *rwClientImpl) PExpire(ctx context.Context, key string, expiration time.Duration) *BoolCmd {
	return c.pool.Master().PExpire(ctx, key, expiration)
}

func (c *rwClientImpl) PExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd {
	return c.pool.Master().PExpireAt(ctx, key, tm)
}

func (c *rwClientImpl) PTTL(ctx context.Context, key string) *DurationCmd {
	return c.pool.Slave().PTTL(ctx, key)
}

func (c *rwClientImpl) RandomKey(ctx context.Context) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Rename(ctx context.Context, key, newkey string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) RenameNX(ctx context.Context, key, newkey string) *BoolCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Restore(ctx context.Context, key string, ttl time.Duration, value string) *StatusCmd {
	return c.pool.Master().Restore(ctx, key, ttl, value)
}

func (c *rwClientImpl) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *StatusCmd {
	return c.pool.Master().RestoreReplace(ctx, key, ttl, value)
}

func (c *rwClientImpl) Sort(ctx context.Context, key string, sort *Sort) *StringSliceCmd {
	return c.pool.Slave().Sort(ctx, key, sort)
}

func (c *rwClientImpl) SortStore(ctx context.Context, key, store string, sort *Sort) *IntCmd {
	return c.pool.Master().SortStore(ctx, key, store, sort)
}

func (c *rwClientImpl) SortInterfaces(ctx context.Context, key string, sort *Sort) *SliceCmd {
	return c.pool.Slave().SortInterfaces(ctx, key, sort)
}

func (c *rwClientImpl) Touch(ctx context.Context, keys ...string) *IntCmd {
	return c.pool.Master().Touch(ctx, keys...)
}

func (c *rwClientImpl) TTL(ctx context.Context, key string) *DurationCmd {
	return c.pool.Slave().TTL(ctx, key)
}

func (c *rwClientImpl) Type(ctx context.Context, key string) *StatusCmd {
	return c.pool.Slave().Type(ctx, key)
}

func (c *rwClientImpl) Append(ctx context.Context, key, value string) *IntCmd {
	return c.pool.Master().Append(ctx, key, value)
}

func (c *rwClientImpl) Decr(ctx context.Context, key string) *IntCmd {
	return c.pool.Master().Decr(ctx, key)
}

func (c *rwClientImpl) DecrBy(ctx context.Context, key string, decrement int64) *IntCmd {
	return c.pool.Master().DecrBy(ctx, key, decrement)
}

func (c *rwClientImpl) Get(ctx context.Context, key string) *StringCmd {
	return c.pool.Slave().Get(ctx, key)
}

func (c *rwClientImpl) GetRange(ctx context.Context, key string, start, end int64) *StringCmd {
	return c.pool.Slave().GetRange(ctx, key, start, end)
}

func (c *rwClientImpl) GetSet(ctx context.Context, key string, value interface{}) *StringCmd {
	return c.pool.Master().GetSet(ctx, key, value)
}

func (c *rwClientImpl) Incr(ctx context.Context, key string) *IntCmd {
	return c.pool.Master().Incr(ctx, key)
}

func (c *rwClientImpl) IncrBy(ctx context.Context, key string, value int64) *IntCmd {
	return c.pool.Master().IncrBy(ctx, key, value)
}

func (c *rwClientImpl) IncrByFloat(ctx context.Context, key string, value float64) *FloatCmd {
	return c.pool.Master().IncrByFloat(ctx, key, value)
}

func (c *rwClientImpl) MGet(ctx context.Context, keys ...string) *SliceCmd {
	return c.pool.Slave().MGet(ctx, keys...)
}

func (c *rwClientImpl) MSet(ctx context.Context, values ...interface{}) *StatusCmd {
	return c.pool.Master().MSet(ctx, values...)
}

func (c *rwClientImpl) MSetNX(ctx context.Context, values ...interface{}) *BoolCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
	return c.pool.Master().Set(ctx, key, value, expiration)
}

func (c *rwClientImpl) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
	return c.pool.Master().SetEX(ctx, key, value, expiration)
}

func (c *rwClientImpl) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd {
	return c.pool.Master().SetNX(ctx, key, value, expiration)
}

func (c *rwClientImpl) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd {
	return c.pool.Master().SetXX(ctx, key, value, expiration)
}

func (c *rwClientImpl) SetRange(ctx context.Context, key string, offset int64, value string) *IntCmd {
	return c.pool.Master().SetRange(ctx, key, offset, value)
}

func (c *rwClientImpl) StrLen(ctx context.Context, key string) *IntCmd {
	return c.pool.Slave().StrLen(ctx, key)
}

func (c *rwClientImpl) GetBit(ctx context.Context, key string, offset int64) *IntCmd {
	return c.pool.Slave().GetBit(ctx, key, offset)
}

func (c *rwClientImpl) SetBit(ctx context.Context, key string, offset int64, value int) *IntCmd {
	return c.pool.Master().SetBit(ctx, key, offset, value)
}

func (c *rwClientImpl) BitCount(ctx context.Context, key string, bitCount *BitCount) *IntCmd {
	return c.pool.Slave().BitCount(ctx, key, bitCount)
}

func (c *rwClientImpl) BitOpAnd(ctx context.Context, destKey string, keys ...string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) BitOpOr(ctx context.Context, destKey string, keys ...string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) BitOpXor(ctx context.Context, destKey string, keys ...string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) BitOpNot(ctx context.Context, destKey string, key string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *IntCmd {
	return c.pool.Slave().BitPos(ctx, key, bit, pos...)
}

func (c *rwClientImpl) BitField(ctx context.Context, key string, args ...interface{}) *IntSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Scan(ctx context.Context, cursor uint64, match string, count int64) *ScanCmd {
	panic("not implemented")
}

func (c *rwClientImpl) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd {
	return c.pool.Slave().SScan(ctx, key, cursor, match, count)
}

func (c *rwClientImpl) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd {
	return c.pool.Slave().HScan(ctx, key, cursor, match, count)
}

func (c *rwClientImpl) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd {
	return c.pool.Slave().ZScan(ctx, key, cursor, match, count)
}

func (c *rwClientImpl) HDel(ctx context.Context, key string, fields ...string) *IntCmd {
	return c.pool.Master().HDel(ctx, key, fields...)
}

func (c *rwClientImpl) HExists(ctx context.Context, key, field string) *BoolCmd {
	return c.pool.Slave().HExists(ctx, key, field)
}

func (c *rwClientImpl) HGet(ctx context.Context, key, field string) *StringCmd {
	return c.pool.Slave().HGet(ctx, key, field)
}

func (c *rwClientImpl) HGetAll(ctx context.Context, key string) *StringStringMapCmd {
	return c.pool.Slave().HGetAll(ctx, key)
}

func (c *rwClientImpl) HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmd {
	return c.pool.Master().HIncrBy(ctx, key, field, incr)
}

func (c *rwClientImpl) HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmd {
	return c.pool.Master().HIncrByFloat(ctx, key, field, incr)
}

func (c *rwClientImpl) HKeys(ctx context.Context, key string) *StringSliceCmd {
	return c.pool.Slave().HKeys(ctx, key)
}

func (c *rwClientImpl) HLen(ctx context.Context, key string) *IntCmd {
	return c.pool.Slave().HLen(ctx, key)
}

func (c *rwClientImpl) HMGet(ctx context.Context, key string, fields ...string) *SliceCmd {
	return c.pool.Slave().HMGet(ctx, key, fields...)
}

func (c *rwClientImpl) HSet(ctx context.Context, key string, values ...interface{}) *IntCmd {
	return c.pool.Master().HSet(ctx, key, values...)
}

func (c *rwClientImpl) HMSet(ctx context.Context, key string, values ...interface{}) *BoolCmd {
	return c.pool.Master().HMSet(ctx, key, values...)
}

func (c *rwClientImpl) HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmd {
	return c.pool.Master().HSetNX(ctx, key, field, value)
}

func (c *rwClientImpl) HVals(ctx context.Context, key string) *StringSliceCmd {
	return c.pool.Slave().HVals(ctx, key)
}

func (c *rwClientImpl) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) LIndex(ctx context.Context, key string, index int64) *StringCmd {
	return c.pool.Slave().LIndex(ctx, key, index)
}

func (c *rwClientImpl) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *IntCmd {
	return c.pool.Master().LInsert(ctx, key, op, pivot, value)
}

func (c *rwClientImpl) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *IntCmd {
	return c.pool.Master().LInsertBefore(ctx, key, pivot, value)
}

func (c *rwClientImpl) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *IntCmd {
	return c.pool.Master().LInsertAfter(ctx, key, pivot, value)
}

func (c *rwClientImpl) LLen(ctx context.Context, key string) *IntCmd {
	return c.pool.Slave().LLen(ctx, key)
}

func (c *rwClientImpl) LPos(ctx context.Context, key string, value string, args redis.LPosArgs) *IntCmd {
	return c.pool.Slave().LPos(ctx, key, value, args)
}

func (c *rwClientImpl) LPosCount(ctx context.Context, key string, value string, count int64, args redis.LPosArgs) *IntSliceCmd {
	return c.pool.Slave().LPosCount(ctx, key, value, count, args)
}

func (c *rwClientImpl) LPop(ctx context.Context, key string) *StringCmd {
	return c.pool.Master().LPop(ctx, key)
}

func (c *rwClientImpl) LPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	return c.pool.Master().LPush(ctx, key, values...)
}

func (c *rwClientImpl) LPushX(ctx context.Context, key string, values ...interface{}) *IntCmd {
	return c.pool.Master().LPushX(ctx, key, values...)
}

func (c *rwClientImpl) LRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	return c.pool.Slave().LRange(ctx, key, start, stop)
}

func (c *rwClientImpl) LRem(ctx context.Context, key string, count int64, value interface{}) *IntCmd {
	return c.pool.Master().LRem(ctx, key, count, value)
}

func (c *rwClientImpl) LSet(ctx context.Context, key string, index int64, value interface{}) *StatusCmd {
	return c.pool.Master().LSet(ctx, key, index, value)
}

func (c *rwClientImpl) LTrim(ctx context.Context, key string, start, stop int64) *StatusCmd {
	return c.pool.Master().LTrim(ctx, key, start, stop)
}

func (c *rwClientImpl) RPop(ctx context.Context, key string) *StringCmd {
	return c.pool.Master().RPop(ctx, key)
}

func (c *rwClientImpl) RPopLPush(ctx context.Context, source, destination string) *StringCmd {
	return c.pool.Master().RPopLPush(ctx, source, destination)
}

func (c *rwClientImpl) RPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	return c.pool.Master().RPush(ctx, key, values...)
}

func (c *rwClientImpl) RPushX(ctx context.Context, key string, values ...interface{}) *IntCmd {
	return c.pool.Master().RPushX(ctx, key, values...)
}

func (c *rwClientImpl) SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd {
	return c.pool.Master().SAdd(ctx, key, members...)
}

func (c *rwClientImpl) SCard(ctx context.Context, key string) *IntCmd {
	return c.pool.Slave().SCard(ctx, key)
}

func (c *rwClientImpl) SDiff(ctx context.Context, keys ...string) *StringSliceCmd {
	return c.pool.Slave().SDiff(ctx, keys...)
}

func (c *rwClientImpl) SDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd {
	return c.pool.Master().SDiffStore(ctx, destination, keys...)
}

func (c *rwClientImpl) SInter(ctx context.Context, keys ...string) *StringSliceCmd {
	return c.pool.Slave().SInter(ctx, keys...)
}

func (c *rwClientImpl) SInterStore(ctx context.Context, destination string, keys ...string) *IntCmd {
	return c.pool.Master().SInterStore(ctx, destination, keys...)
}

func (c *rwClientImpl) SIsMember(ctx context.Context, key string, member interface{}) *BoolCmd {
	return c.pool.Slave().SIsMember(ctx, key, member)
}

func (c *rwClientImpl) SMembers(ctx context.Context, key string) *StringSliceCmd {
	return c.pool.Slave().SMembers(ctx, key)
}

func (c *rwClientImpl) SMembersMap(ctx context.Context, key string) *StringStructMapCmd {
	return c.pool.Slave().SMembersMap(ctx, key)
}

func (c *rwClientImpl) SMove(ctx context.Context, source, destination string, member interface{}) *BoolCmd {
	return c.pool.Master().SMove(ctx, source, destination, member)
}

func (c *rwClientImpl) SPop(ctx context.Context, key string) *StringCmd {
	return c.pool.Master().SPop(ctx, key)
}

func (c *rwClientImpl) SPopN(ctx context.Context, key string, count int64) *StringSliceCmd {
	return c.pool.Master().SPopN(ctx, key, count)
}

func (c *rwClientImpl) SRandMember(ctx context.Context, key string) *StringCmd {
	return c.pool.Slave().SRandMember(ctx, key)
}

func (c *rwClientImpl) SRandMemberN(ctx context.Context, key string, count int64) *StringSliceCmd {
	return c.pool.Slave().SRandMemberN(ctx, key, count)
}

func (c *rwClientImpl) SRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	return c.pool.Master().SRem(ctx, key, members)
}

func (c *rwClientImpl) SUnion(ctx context.Context, keys ...string) *StringSliceCmd {
	return c.pool.Slave().SUnion(ctx, keys...)
}

func (c *rwClientImpl) SUnionStore(ctx context.Context, destination string, keys ...string) *IntCmd {
	return c.pool.Master().SUnionStore(ctx, destination, keys...)
}

func (c *rwClientImpl) XAdd(ctx context.Context, a *XAddArgs) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XDel(ctx context.Context, stream string, ids ...string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XLen(ctx context.Context, stream string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XRange(ctx context.Context, stream, start, stop string) *XMessageSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XRangeN(ctx context.Context, stream, start, stop string, count int64) *XMessageSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XRevRange(ctx context.Context, stream string, start, stop string) *XMessageSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) *XMessageSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XRead(ctx context.Context, a *XReadArgs) *XStreamSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XReadStreams(ctx context.Context, streams ...string) *XStreamSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XGroupCreate(ctx context.Context, stream, group, start string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XGroupSetID(ctx context.Context, stream, group, start string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XGroupDestroy(ctx context.Context, stream, group string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XReadGroup(ctx context.Context, a *XReadGroupArgs) *XStreamSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XAck(ctx context.Context, stream, group string, ids ...string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XPending(ctx context.Context, stream, group string) *XPendingCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XPendingExt(ctx context.Context, a *XPendingExtArgs) *XPendingExtCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XClaim(ctx context.Context, a *XClaimArgs) *XMessageSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XClaimJustID(ctx context.Context, a *XClaimArgs) *StringSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XTrim(ctx context.Context, key string, maxLen int64) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XTrimApprox(ctx context.Context, key string, maxLen int64) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XInfoGroups(ctx context.Context, key string) *XInfoGroupsCmd {
	panic("not implemented")
}

func (c *rwClientImpl) XInfoStream(ctx context.Context, key string) *XInfoStreamCmd {
	panic("not implemented")
}

func (c *rwClientImpl) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd {
	panic("not implemented")
}

func (c *rwClientImpl) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ZAdd(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.pool.Master().ZAdd(ctx, key, members...)
}

func (c *rwClientImpl) ZAddNX(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.pool.Master().ZAddNX(ctx, key, members...)
}

func (c *rwClientImpl) ZAddXX(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.pool.Master().ZAddXX(ctx, key, members...)
}

func (c *rwClientImpl) ZAddCh(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.pool.Master().ZAddCh(ctx, key, members...)
}

func (c *rwClientImpl) ZAddNXCh(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.pool.Master().ZAddNXCh(ctx, key, members...)
}

func (c *rwClientImpl) ZAddXXCh(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.pool.Master().ZAddXXCh(ctx, key, members...)
}

func (c *rwClientImpl) ZIncr(ctx context.Context, key string, member *Z) *FloatCmd {
	return c.pool.Master().ZIncr(ctx, key, member)
}

func (c *rwClientImpl) ZIncrNX(ctx context.Context, key string, member *Z) *FloatCmd {
	return c.pool.Master().ZIncrNX(ctx, key, member)
}

func (c *rwClientImpl) ZIncrXX(ctx context.Context, key string, member *Z) *FloatCmd {
	return c.pool.Master().ZIncrXX(ctx, key, member)
}

func (c *rwClientImpl) ZCard(ctx context.Context, key string) *IntCmd {
	return c.pool.Slave().ZCard(ctx, key)
}

func (c *rwClientImpl) ZCount(ctx context.Context, key, min, max string) *IntCmd {
	return c.pool.Slave().ZCount(ctx, key, min, max)
}

func (c *rwClientImpl) ZLexCount(ctx context.Context, key, min, max string) *IntCmd {
	return c.pool.Slave().ZLexCount(ctx, key, min, max)
}

func (c *rwClientImpl) ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd {
	return c.pool.Master().ZIncrBy(ctx, key, increment, member)
}

func (c *rwClientImpl) ZInterStore(ctx context.Context, destination string, store *ZStore) *IntCmd {
	return c.pool.Master().ZInterStore(ctx, destination, store)
}

func (c *rwClientImpl) ZPopMax(ctx context.Context, key string, count ...int64) *ZSliceCmd {
	return c.pool.Master().ZPopMax(ctx, key, count...)
}

func (c *rwClientImpl) ZPopMin(ctx context.Context, key string, count ...int64) *ZSliceCmd {
	return c.pool.Master().ZPopMin(ctx, key, count...)
}

func (c *rwClientImpl) ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	return c.pool.Slave().ZRange(ctx, key, start, stop)
}

func (c *rwClientImpl) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd {
	return c.pool.Slave().ZRangeWithScores(ctx, key, start, stop)
}

func (c *rwClientImpl) ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.pool.Slave().ZRangeByScore(ctx, key, opt)
}

func (c *rwClientImpl) ZRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.pool.Slave().ZRangeByLex(ctx, key, opt)
}

func (c *rwClientImpl) ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd {
	return c.pool.Slave().ZRangeByScoreWithScores(ctx, key, opt)
}

func (c *rwClientImpl) ZRank(ctx context.Context, key, member string) *IntCmd {
	return c.pool.Slave().ZRank(ctx, key, member)
}

func (c *rwClientImpl) ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	return c.pool.Master().ZRem(ctx, key, members...)
}

func (c *rwClientImpl) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd {
	return c.pool.Master().ZRemRangeByRank(ctx, key, start, stop)
}

func (c *rwClientImpl) ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd {
	return c.pool.Master().ZRemRangeByScore(ctx, key, min, max)
}

func (c *rwClientImpl) ZRemRangeByLex(ctx context.Context, key, min, max string) *IntCmd {
	return c.pool.Master().ZRemRangeByLex(ctx, key, min, max)
}

func (c *rwClientImpl) ZRevRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	return c.pool.Slave().ZRevRange(ctx, key, start, stop)
}

func (c *rwClientImpl) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd {
	return c.pool.Slave().ZRevRangeWithScores(ctx, key, start, stop)
}

func (c *rwClientImpl) ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.pool.Slave().ZRevRangeByScore(ctx, key, opt)
}

func (c *rwClientImpl) ZRevRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.pool.Slave().ZRevRangeByLex(ctx, key, opt)
}

func (c *rwClientImpl) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd {
	return c.pool.Slave().ZRevRangeByScoreWithScores(ctx, key, opt)
}

func (c *rwClientImpl) ZRevRank(ctx context.Context, key, member string) *IntCmd {
	return c.pool.Slave().ZRevRank(ctx, key, member)
}

func (c *rwClientImpl) ZScore(ctx context.Context, key, member string) *FloatCmd {
	return c.pool.Slave().ZScore(ctx, key, member)
}

func (c *rwClientImpl) ZUnionStore(ctx context.Context, dest string, store *ZStore) *IntCmd {
	return c.pool.Master().ZUnionStore(ctx, dest, store)
}

func (c *rwClientImpl) PFAdd(ctx context.Context, key string, els ...interface{}) *IntCmd {
	return c.pool.Master().PFAdd(ctx, key, els...)
}

func (c *rwClientImpl) PFCount(ctx context.Context, keys ...string) *IntCmd {
	return c.pool.Slave().PFCount(ctx, keys...)
}

func (c *rwClientImpl) PFMerge(ctx context.Context, dest string, keys ...string) *StatusCmd {
	return c.pool.Master().PFMerge(ctx, dest, keys...)
}

func (c *rwClientImpl) BgRewriteAOF(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) BgSave(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClientKill(ctx context.Context, ipPort string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClientKillByFilter(ctx context.Context, keys ...string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClientList(ctx context.Context) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClientPause(ctx context.Context, dur time.Duration) *BoolCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClientID(ctx context.Context) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ConfigGet(ctx context.Context, parameter string) *SliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ConfigResetStat(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ConfigSet(ctx context.Context, parameter, value string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ConfigRewrite(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) DBSize(ctx context.Context) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) FlushAll(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) FlushAllAsync(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) FlushDB(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) FlushDBAsync(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Info(ctx context.Context, section ...string) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) LastSave(ctx context.Context) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Save(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Shutdown(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ShutdownSave(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ShutdownNoSave(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) SlaveOf(ctx context.Context, host, port string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Time(ctx context.Context) *TimeCmd {
	panic("not implemented")
}

func (c *rwClientImpl) DebugObject(ctx context.Context, key string) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ReadOnly(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ReadWrite(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) MemoryUsage(ctx context.Context, key string, samples ...int) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *Cmd {
	panic("not implemented")
}

func (c *rwClientImpl) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *Cmd {
	panic("not implemented")
}

func (c *rwClientImpl) ScriptExists(ctx context.Context, hashes ...string) *BoolSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ScriptFlush(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ScriptKill(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ScriptLoad(ctx context.Context, script string) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) Publish(ctx context.Context, channel string, message interface{}) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) PubSubChannels(ctx context.Context, pattern string) *StringSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) PubSubNumSub(ctx context.Context, channels ...string) *StringIntMapCmd {
	panic("not implemented")
}

func (c *rwClientImpl) PubSubNumPat(ctx context.Context) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterSlots(ctx context.Context) *ClusterSlotsCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterNodes(ctx context.Context) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterMeet(ctx context.Context, host, port string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterForget(ctx context.Context, nodeID string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterReplicate(ctx context.Context, nodeID string) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterResetSoft(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterResetHard(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterInfo(ctx context.Context) *StringCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterKeySlot(ctx context.Context, key string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *StringSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterCountFailureReports(ctx context.Context, nodeID string) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterCountKeysInSlot(ctx context.Context, slot int) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterDelSlots(ctx context.Context, slots ...int) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterDelSlotsRange(ctx context.Context, min, max int) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterSaveConfig(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterSlaves(ctx context.Context, nodeID string) *StringSliceCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterFailover(ctx context.Context) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterAddSlots(ctx context.Context, slots ...int) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) ClusterAddSlotsRange(ctx context.Context, min, max int) *StatusCmd {
	panic("not implemented")
}

func (c *rwClientImpl) GeoAdd(ctx context.Context, key string, geoLocation ...*GeoLocation) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) GeoPos(ctx context.Context, key string, members ...string) *GeoPosCmd {
	panic("not implemented")
}

func (c *rwClientImpl) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery) *GeoLocationCmd {
	panic("not implemented")
}

func (c *rwClientImpl) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) GeoRadiusByMember(ctx context.Context, key, member string, query *GeoRadiusQuery) *GeoLocationCmd {
	panic("not implemented")
}

func (c *rwClientImpl) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *GeoRadiusQuery) *IntCmd {
	panic("not implemented")
}

func (c *rwClientImpl) GeoDist(ctx context.Context, key string, member1, member2, unit string) *FloatCmd {
	panic("not implemented")
}

func (c *rwClientImpl) GeoHash(ctx context.Context, key string, members ...string) *StringSliceCmd {
	panic("not implemented")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) GetEx(ctx context.Context, key string, expiration time.Duration) *redis.StringCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) GetDel(ctx context.Context, key string) *redis.StringCmd {
	// TODO implement me
	panic("not supported yet")
}

func (c *rwClientImpl) SetArgs(ctx context.Context, key string, value interface{}, a SetArgs) *redis.StatusCmd {
	return c.pool.Master().SetArgs(ctx, key, value, a)
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) HRandField(ctx context.Context, key string, count int, withValues bool) *redis.StringSliceCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) LPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) LMove(ctx context.Context, source, destination, srcpos, destpos string) *redis.StringCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) SMIsMember(ctx context.Context, key string, members ...interface{}) *redis.BoolSliceCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) XInfoConsumers(ctx context.Context, key string, group string) *redis.XInfoConsumersCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) ZInter(ctx context.Context, store *redis.ZStore) *redis.StringSliceCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) ZInterWithScores(ctx context.Context, store *redis.ZStore) *redis.ZSliceCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) ZMScore(ctx context.Context, key string, members ...string) *redis.FloatSliceCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) ZRandMember(ctx context.Context, key string, count int, withScores bool) *redis.StringSliceCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) ZDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) ZDiffWithScores(ctx context.Context, keys ...string) *redis.ZSliceCmd {
	// TODO implement me
	panic("not supported yet")
}

// NOT SUPPORTED YET, DO NOT USE
func (c *rwClientImpl) ZDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	// TODO implement me
	panic("not supported yet")
}

func (c *rwClientImpl) XTrimMaxLen(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	panic("not supported yet")
}
func (c *rwClientImpl) XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) *redis.IntCmd {
	panic("not supported yet")
}
func (c *rwClientImpl) XTrimMinID(ctx context.Context, key string, minID string) *redis.IntCmd {
	panic("not supported yet")
}
func (c *rwClientImpl) XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) *redis.IntCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) XGroupCreateConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ZAddArgs(ctx context.Context, key string, args redis.ZAddArgs) *redis.IntCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ZAddArgsIncr(ctx context.Context, key string, args redis.ZAddArgs) *redis.FloatCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ZRangeArgs(ctx context.Context, z redis.ZRangeArgs) *redis.StringSliceCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ZRangeArgsWithScores(ctx context.Context, z redis.ZRangeArgs) *redis.ZSliceCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ZRangeStore(ctx context.Context, dst string, z redis.ZRangeArgs) *redis.IntCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ZUnion(ctx context.Context, store redis.ZStore) *redis.StringSliceCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ZUnionWithScores(ctx context.Context, store redis.ZStore) *redis.ZSliceCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) *redis.StringCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) Copy(ctx context.Context, sourceKey string, destKey string, db int, replace bool) *redis.IntCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ExpireGT(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ExpireLT(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ExpireNX(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) ExpireXX(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) *redis.StringSliceCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) *redis.GeoSearchLocationCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) GeoSearchStore(ctx context.Context, key, store string, q *redis.GeoSearchStoreQuery) *redis.IntCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) RPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) XAutoClaim(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) XAutoClaimJustID(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimJustIDCmd {
	panic("not supported yet")
}

func (c *rwClientImpl) XInfoStreamFull(ctx context.Context, key string, count int) *redis.XInfoStreamFullCmd {
	panic("not supported yet")
}
