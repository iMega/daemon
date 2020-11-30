package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

type faker struct{}

// ErrFaker is a faker
const ErrFaker = "it is a faker database"

func fakeIntCmd() *redis.IntCmd {
	return redis.NewIntResult(0, errors.New(ErrFaker))
}

func fakeStringCmd() *redis.StringCmd {
	return redis.NewStringResult("", errors.New(ErrFaker))
}

func fakeBoolCmd() *redis.BoolCmd {
	return redis.NewBoolResult(false, errors.New(ErrFaker))
}

func fakeDurationCmd() *redis.DurationCmd {
	return redis.NewDurationResult(0, errors.New(ErrFaker))
}

func fakeStatusCmd() *redis.StatusCmd {
	return redis.NewStatusResult("", errors.New(ErrFaker))
}

func fakeStringSliceCmd() *redis.StringSliceCmd {
	return redis.NewStringSliceResult(nil, errors.New(ErrFaker))
}

func fakeSliceCmd() *redis.SliceCmd {
	return redis.NewSliceResult(nil, errors.New(ErrFaker))
}

func fakeScanCmd() *redis.ScanCmd {
	return redis.NewScanCmdResult(nil, 0, errors.New(ErrFaker))
}

func fakeFloatCmd() *redis.FloatCmd {
	return redis.NewFloatResult(0, errors.New(ErrFaker))
}

func fakeCmdResult() *redis.Cmd {
	return redis.NewCmdResult(nil, errors.New(ErrFaker))
}

func fakeZSliceCmdResult() *redis.ZSliceCmd {
	return redis.NewZSliceCmdResult(nil, errors.New(ErrFaker))
}

func fakeGeoLocationCmdResult() *redis.GeoLocationCmd {
	return redis.NewGeoLocationCmdResult(nil, errors.New(ErrFaker))
}

// Release interface

func (*faker) Watch(fn func(*redis.Tx) error, keys ...string) error {
	return errors.New(ErrFaker)
}

func (*faker) Process(redis.Cmder) error {
	return errors.New(ErrFaker)
}

func (*faker) WrapProcess(fn func(oldProcess func(redis.Cmder) error) func(redis.Cmder) error) {}

func (*faker) WrapProcessPipeline(fn func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error) {
}

func (*faker) Subscribe(channels ...string) *redis.PubSub {
	return nil
}

func (*faker) PSubscribe(channels ...string) *redis.PubSub {
	return nil
}

func (*faker) Close() error {
	return errors.New(ErrFaker)
}

func (*faker) Pipeline() redis.Pipeliner {
	return nil
}

func (*faker) Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return nil, errors.New(ErrFaker)
}

func (*faker) TxPipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return nil, errors.New(ErrFaker)
}

func (*faker) TxPipeline() redis.Pipeliner {
	return nil
}

func (*faker) Command() *redis.CommandsInfoCmd {
	return redis.NewCommandsInfoCmdResult(nil, errors.New(ErrFaker))
}

func (*faker) ClientGetName() *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) Echo(message interface{}) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) Ping() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) Quit() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) Del(keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) Unlink(keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) Dump(key string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) Exists(keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) Keys(pattern string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) Migrate(host, port, key string, db int64, timeout time.Duration) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) Move(key string, db int64) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) ObjectRefCount(key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ObjectEncoding(key string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) ObjectIdleTime(key string) *redis.DurationCmd {
	return fakeDurationCmd()
}

func (*faker) Persist(key string) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) PExpire(key string, expiration time.Duration) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) PExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) PTTL(key string) *redis.DurationCmd {
	return fakeDurationCmd()
}

func (*faker) RandomKey() *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) Rename(key, newkey string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) RenameNX(key, newkey string) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) Restore(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) Sort(key string, sort *redis.Sort) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) SortStore(key, store string, sort *redis.Sort) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) SortInterfaces(key string, sort *redis.Sort) *redis.SliceCmd {
	return fakeSliceCmd()
}

func (*faker) Touch(keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) TTL(key string) *redis.DurationCmd {
	return fakeDurationCmd()
}

func (*faker) Type(key string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return fakeScanCmd()
}

func (*faker) SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return fakeScanCmd()
}

func (*faker) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return fakeScanCmd()
}

func (*faker) ZScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return fakeScanCmd()
}

func (*faker) Append(key, value string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) BitOpNot(destKey string, key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) BitPos(key string, bit int64, pos ...int64) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) Decr(key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) DecrBy(key string, decrement int64) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) Get(key string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) GetBit(key string, offset int64) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) GetRange(key string, start, end int64) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) GetSet(key string, value interface{}) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) Incr(key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) IncrBy(key string, value int64) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) IncrByFloat(key string, value float64) *redis.FloatCmd {
	return fakeFloatCmd()
}

func (*faker) MGet(keys ...string) *redis.SliceCmd {
	return fakeSliceCmd()
}

func (*faker) MSet(pairs ...interface{}) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) MSetNX(pairs ...interface{}) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) SetBit(key string, offset int64, value int) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) SetRange(key string, offset int64, value string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) StrLen(key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) HDel(key string, fields ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) HExists(key, field string) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) HGet(key, field string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) HGetAll(key string) *redis.StringStringMapCmd {
	return redis.NewStringStringMapResult(nil, errors.New(ErrFaker))
}

func (*faker) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	return fakeFloatCmd()
}

func (*faker) HKeys(key string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) HLen(key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) HMGet(key string, fields ...string) *redis.SliceCmd {
	return fakeSliceCmd()
}

func (*faker) HMSet(key string, fields map[string]interface{}) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) HSet(key, field string, value interface{}) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) HVals(key string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) LIndex(key string, index int64) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) LInsert(key, op string, pivot, value interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) LInsertBefore(key string, pivot, value interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) LInsertAfter(key string, pivot, value interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) LLen(key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) LPop(key string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) LPush(key string, values ...interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) LPushX(key string, value interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) LTrim(key string, start, stop int64) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) RPop(key string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) RPopLPush(source, destination string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) RPush(key string, values ...interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) RPushX(key string, value interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) SCard(key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) SDiff(keys ...string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) SInter(keys ...string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) SInterStore(destination string, keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) SIsMember(key string, member interface{}) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) SMembers(key string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) SMembersMap(key string) *redis.StringStructMapCmd {
	// TODO: not implement redis.NewStringStructMapResult in tests
	return nil
}

func (*faker) SMove(source, destination string, member interface{}) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) SPop(key string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) SPopN(key string, count int64) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) SRandMember(key string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) SRem(key string, members ...interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) SUnion(keys ...string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) SUnionStore(destination string, keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) XAdd(a *redis.XAddArgs) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) XDel(stream string, ids ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) XLen(stream string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) XRange(stream, start, stop string) *redis.XMessageSliceCmd {
	// TODO: not implement redis.NewXMessageSliceResult in tests
	return nil
}

func (*faker) XRangeN(stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	// TODO: not implement redis.NewXMessageSliceResult in tests
	return nil
}

func (*faker) XRevRange(stream string, start, stop string) *redis.XMessageSliceCmd {
	// TODO: not implement redis.NewXMessageSliceResult in tests
	return nil
}

func (*faker) XRevRangeN(stream string, start, stop string, count int64) *redis.XMessageSliceCmd {
	// TODO: not implement redis.NewXMessageSliceResult in tests
	return nil
}

func (*faker) XRead(a *redis.XReadArgs) *redis.XStreamSliceCmd {
	// TODO: not implement redis.NewXStreamSliceResult in tests
	return nil
}

func (*faker) XReadStreams(streams ...string) *redis.XStreamSliceCmd {
	// TODO: not implement redis.NewXStreamSliceResult in tests
	return nil
}

func (*faker) XGroupCreate(stream, group, start string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) XGroupCreateMkStream(stream, group, start string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) XGroupSetID(stream, group, start string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) XGroupDestroy(stream, group string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) XGroupDelConsumer(stream, group, consumer string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) XReadGroup(a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	// TODO: not implement redis.NewXStreamSliceResult in tests
	return nil
}

func (*faker) XAck(stream, group string, ids ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) XPending(stream, group string) *redis.XPendingCmd {
	// TODO: not implement redis.NewXPendingResult in tests
	return nil
}

func (*faker) XPendingExt(a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	// TODO: not implement redis.NewXPendingExtResult in tests
	return nil
}

func (*faker) XClaim(a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	// TODO: not implement redis.NewXMessageSliceResult in tests
	return nil
}

func (*faker) XClaimJustID(a *redis.XClaimArgs) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) XTrim(key string, maxLen int64) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) XTrimApprox(key string, maxLen int64) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) BZPopMax(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	// TODO: not implement redis.NewZWithKeyResult in tests
	return nil
}

func (*faker) BZPopMin(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	// TODO: not implement redis.NewZWithKeyResult in tests
	return nil
}

func (*faker) ZAdd(key string, members ...redis.Z) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZAddNX(key string, members ...redis.Z) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZAddXX(key string, members ...redis.Z) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZAddCh(key string, members ...redis.Z) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZAddNXCh(key string, members ...redis.Z) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZAddXXCh(key string, members ...redis.Z) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZIncr(key string, member redis.Z) *redis.FloatCmd {
	return fakeFloatCmd()
}

func (*faker) ZIncrNX(key string, member redis.Z) *redis.FloatCmd {
	return fakeFloatCmd()
}

func (*faker) ZIncrXX(key string, member redis.Z) *redis.FloatCmd {
	return fakeFloatCmd()
}

func (*faker) ZCard(key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZCount(key, min, max string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZLexCount(key, min, max string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	return fakeFloatCmd()
}

func (*faker) ZInterStore(destination string, store redis.ZStore, keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZPopMax(key string, count ...int64) *redis.ZSliceCmd {
	return fakeZSliceCmdResult()
}

func (*faker) ZPopMin(key string, count ...int64) *redis.ZSliceCmd {
	return fakeZSliceCmdResult()
}

func (*faker) ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return fakeZSliceCmdResult()
}

func (*faker) ZRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) ZRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	return fakeZSliceCmdResult()
}

func (*faker) ZRank(key, member string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZRem(key string, members ...interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZRemRangeByScore(key, min, max string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZRemRangeByLex(key, min, max string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZRevRange(key string, start, stop int64) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) ZRevRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return fakeZSliceCmdResult()
}

func (*faker) ZRevRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) ZRevRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	return fakeZSliceCmdResult()
}

func (*faker) ZRevRank(key, member string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ZScore(key, member string) *redis.FloatCmd {
	return fakeFloatCmd()
}

func (*faker) ZUnionStore(dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) PFAdd(key string, els ...interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) PFCount(keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) BgRewriteAOF() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) BgSave() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClientKill(ipPort string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClientKillByFilter(keys ...string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ClientList() *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) ClientPause(dur time.Duration) *redis.BoolCmd {
	return fakeBoolCmd()
}

func (*faker) ClientID() *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ConfigGet(parameter string) *redis.SliceCmd {
	return fakeSliceCmd()
}

func (*faker) ConfigResetStat() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ConfigSet(parameter, value string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ConfigRewrite() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) DBSize() *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) FlushAll() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) FlushAllAsync() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) FlushDB() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) FlushDBAsync() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) Info(section ...string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) LastSave() *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) Save() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) Shutdown() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ShutdownSave() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ShutdownNoSave() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) SlaveOf(host, port string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) Time() *redis.TimeCmd {
	// TODO: not implement redis.NewTimeResult in tests
	return nil
}

func (*faker) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	return fakeCmdResult()
}

func (*faker) EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return fakeCmdResult()
}

func (*faker) ScriptExists(hashes ...string) *redis.BoolSliceCmd {
	return redis.NewBoolSliceResult(nil, errors.New(ErrFaker))
}

func (*faker) ScriptFlush() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ScriptKill() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ScriptLoad(script string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) DebugObject(key string) *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) Publish(channel string, message interface{}) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) PubSubChannels(pattern string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) PubSubNumSub(channels ...string) *redis.StringIntMapCmd {
	return redis.NewStringIntMapCmdResult(nil, errors.New(ErrFaker))
}

func (*faker) PubSubNumPat() *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ClusterSlots() *redis.ClusterSlotsCmd {
	return redis.NewClusterSlotsCmdResult(nil, errors.New(ErrFaker))
}

func (*faker) ClusterNodes() *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) ClusterMeet(host, port string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterForget(nodeID string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterReplicate(nodeID string) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterResetSoft() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterResetHard() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterInfo() *redis.StringCmd {
	return fakeStringCmd()
}

func (*faker) ClusterKeySlot(key string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ClusterGetKeysInSlot(slot int, count int) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) ClusterCountFailureReports(nodeID string) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ClusterCountKeysInSlot(slot int) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) ClusterDelSlots(slots ...int) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterDelSlotsRange(min, max int) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterSaveConfig() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterSlaves(nodeID string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) ClusterFailover() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterAddSlots(slots ...int) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ClusterAddSlotsRange(min, max int) *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	return fakeIntCmd()
}

func (*faker) GeoPos(key string, members ...string) *redis.GeoPosCmd {
	// TODO: not implement redis.NewGeoPosResult in tests
	return nil
}

func (*faker) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return fakeGeoLocationCmdResult()
}

func (*faker) GeoRadiusRO(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return fakeGeoLocationCmdResult()
}

func (*faker) GeoRadiusByMember(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return fakeGeoLocationCmdResult()
}

func (*faker) GeoRadiusByMemberRO(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return fakeGeoLocationCmdResult()
}

func (*faker) GeoDist(key string, member1, member2, unit string) *redis.FloatCmd {
	return fakeFloatCmd()
}

func (*faker) GeoHash(key string, members ...string) *redis.StringSliceCmd {
	return fakeStringSliceCmd()
}

func (*faker) ReadOnly() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) ReadWrite() *redis.StatusCmd {
	return fakeStatusCmd()
}

func (*faker) MemoryUsage(key string, samples ...int) *redis.IntCmd {
	return fakeIntCmd()
}
