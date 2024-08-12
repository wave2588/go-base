package redis

import "github.com/go-redis/redis/v8"

// redis.go
const Nil = redis.Nil

type Client = redis.Cmdable
type Hook = redis.Hook
type PoolStats = redis.PoolStats
type Conn = redis.Conn

// pipeline.go
type Pipeline = redis.Pipeline
type Pipeliner = redis.Pipeliner

// commands.go
const KeepTTL = redis.KeepTTL

type SetArgs = redis.SetArgs

type Cmder = redis.Cmder
type Cmd = redis.Cmd
type SliceCmd = redis.SliceCmd
type StatusCmd = redis.StatusCmd
type IntCmd = redis.IntCmd
type IntSliceCmd = redis.IntSliceCmd
type DurationCmd = redis.DurationCmd
type TimeCmd = redis.TimeCmd
type BoolCmd = redis.BoolCmd
type StringCmd = redis.StringCmd
type FloatCmd = redis.FloatCmd
type StringSliceCmd = redis.StringSliceCmd
type BoolSliceCmd = redis.BoolSliceCmd
type StringStringMapCmd = redis.StringStringMapCmd
type StringIntMapCmd = redis.StringIntMapCmd
type StringStructMapCmd = redis.StringStructMapCmd
type XMessageSliceCmd = redis.XMessageSliceCmd
type XStreamSliceCmd = redis.XStreamSliceCmd
type XPendingCmd = redis.XPendingCmd
type XPendingExtCmd = redis.XPendingExtCmd
type XInfoGroupsCmd = redis.XInfoGroupsCmd
type XInfoStreamCmd = redis.XInfoStreamCmd
type ZSliceCmd = redis.ZSliceCmd
type ZWithKeyCmd = redis.ZWithKeyCmd
type ScanCmd = redis.ScanCmd
type ClusterSlotsCmd = redis.ClusterSlotsCmd
type GeoLocationCmd = redis.GeoLocationCmd
type GeoPosCmd = redis.GeoPosCmd
type CommandsInfoCmd = redis.CommandsInfoCmd
type SlowLogCmd = redis.SlowLogCmd

type Sort = redis.Sort
type BitCount = redis.BitCount
type XAddArgs = redis.XAddArgs
type XReadArgs = redis.XReadArgs
type XReadGroupArgs = redis.XReadGroupArgs
type XPendingExtArgs = redis.XPendingExtArgs
type XClaimArgs = redis.XClaimArgs
type Z = redis.Z
type ZStore = redis.ZStore
type ZRangeBy = redis.ZRangeBy
type GeoLocation = redis.GeoLocation
type GeoRadiusQuery = redis.GeoRadiusQuery

// pubsub.go
type PubSub = redis.PubSub
type Subscription = redis.Subscription
type Message = redis.Message
type Pong = redis.Pong

// tx.go
const TxFailedErr = redis.TxFailedErr

type Tx = redis.Tx
