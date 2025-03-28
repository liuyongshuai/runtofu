// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        redis_client.go
// @date        2025-03-25 10:55

package negoutils

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

func InitRedisClient(hosts string, timeout time.Duration) (ret *RedisClient, err error) {
	ret = &RedisClient{}
	ret.SetHost(hosts).
		SetMaxActive(256).
		SetMaxIdle(8).
		SetIdleTimeout(timeout).
		SetDialConnTimeout(timeout).
		SetDialReadTimeout(timeout).
		SetDialWriteTimeout(timeout)
	err = ret.Init()
	return ret, err
}

type RedisClient struct {
	RedisPool        *redis.Pool
	Hosts            string
	Password         string
	MaxActive        int
	MaxIdle          int
	IdleTimeout      time.Duration
	DialConnTimeout  time.Duration
	DialReadTimeout  time.Duration
	DialWriteTimeout time.Duration
}

func (rc *RedisClient) SetHost(hosts string) *RedisClient {
	rc.Hosts = hosts
	return rc
}

func (rc *RedisClient) SetPassword(password string) *RedisClient {
	rc.Password = password
	return rc
}
func (rc *RedisClient) SetMaxActive(maxActive int) *RedisClient {
	rc.MaxActive = maxActive
	return rc
}
func (rc *RedisClient) SetMaxIdle(maxIdle int) *RedisClient {
	rc.MaxIdle = maxIdle
	return rc
}
func (rc *RedisClient) SetIdleTimeout(idleTimeout time.Duration) *RedisClient {
	rc.IdleTimeout = idleTimeout
	return rc
}
func (rc *RedisClient) SetDialConnTimeout(connTimeout time.Duration) *RedisClient {
	rc.DialConnTimeout = connTimeout
	return rc
}
func (rc *RedisClient) SetDialReadTimeout(readTimeout time.Duration) *RedisClient {
	rc.DialReadTimeout = readTimeout
	return rc
}
func (rc *RedisClient) SetDialWriteTimeout(WriteTimeout time.Duration) *RedisClient {
	rc.DialWriteTimeout = WriteTimeout
	return rc
}

func (rc *RedisClient) Init() error {
	rc.RedisPool = &redis.Pool{
		Wait:        true,
		MaxIdle:     rc.MaxIdle,
		MaxActive:   rc.MaxActive,
		IdleTimeout: rc.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				rc.Hosts,
				redis.DialConnectTimeout(rc.DialConnTimeout),
				redis.DialReadTimeout(rc.DialReadTimeout),
				redis.DialWriteTimeout(rc.DialWriteTimeout),
			)
			if err != nil {
				return nil, err
			}
			if rc.Password != "" {
				err := c.Send("auth", rc.Password)
				if err != nil {
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

// GET命令
func (rc *RedisClient) Get(key string) (val string, err error) {
	val, err = redis.String(rc.Do("GET", key))
	return
}

// MGET命令
func (rc *RedisClient) MGet(keys []string) (vals []string, err error) {
	arr := make([]interface{}, len(keys))
	for k := range arr {
		arr[k] = keys[k]
	}
	vals, err = redis.Strings(rc.Do("MGET", arr...))
	return
}

// DEL命令
func (rc *RedisClient) Del(key string) (ret bool, err error) {
	return redis.Bool(rc.Do("DEL", key))
}

// Expire命令
func (rc *RedisClient) Expire(key string, exptime int64) (ret bool, err error) {
	return redis.Bool(rc.Do("EXPIRE", key, exptime))
}

// SET命令
func (rc *RedisClient) Set(key string, val string) (ret bool, err error) {
	retstr, err := redis.String(rc.Do("SET", key, val))
	if err != nil {
		return false, err
	}
	if retstr == "OK" {
		return true, nil
	}
	return false, nil
}

// SETEX命令
func (rc *RedisClient) SetEX(key string, val string, exptime int64) (ret bool, err error) {
	retstr, err := redis.String(rc.Do("SETEX", key, exptime, val))
	if err != nil {
		return false, err
	}
	if retstr == "OK" {
		return true, nil
	}
	return false, nil
}

// SETNX命令
func (rc *RedisClient) SetNX(key string, val string) (ret bool, err error) {
	return redis.Bool(rc.Do("SETNX", key, val))
}

// EXISTS命令
func (rc *RedisClient) Exists(key string) (ret bool, err error) {
	return redis.Bool(rc.Do("EXISTS", key))
}

// MSETNX命令
func (rc *RedisClient) MSetNX(kvs map[string]string) (err error) {
	size := 2 * len(kvs)
	arr := make([]interface{}, 0, size)
	for k, v := range kvs {
		arr = append(arr, k)
		arr = append(arr, v)
	}
	_, err = rc.Do("MSETNX", arr...)
	return
}

// MSET命令
func (rc *RedisClient) MSet(kvs map[string]string) (err error) {
	size := 2 * len(kvs)
	arr := make([]interface{}, 0, size)
	for k, v := range kvs {
		arr = append(arr, k)
		arr = append(arr, v)
	}
	_, err = rc.Do("MSET", arr...)
	return
}

// SADD命令
func (rc *RedisClient) SAdd(key string, vals ...interface{}) (ret int, err error) {
	arr := make([]interface{}, 0, len(vals)+1)
	arr = append(arr, key)
	arr = append(arr, vals...)
	return redis.Int(rc.Do("SADD", arr...))
}

// SREM命令
func (rc *RedisClient) SRem(key string, vals ...interface{}) (ret int, err error) {
	arr := make([]interface{}, 0, len(vals)+1)
	arr = append(arr, key)
	arr = append(arr, vals...)
	return redis.Int(rc.Do("SREM", arr...))
}

// SMEMBERS命令
func (rc *RedisClient) SMembers(key string) (vals []string, err error) {
	return redis.Strings(rc.Do("SMEMBERS", key))
}

// INCRBY命令
func (rc *RedisClient) IncrBy(key string, val int64) (newVal int64, err error) {
	return redis.Int64(rc.Do("INCRBY", key, val))
}

// ZRANGE命令
func (rc *RedisClient) ZRange(key string, startIndex, count int64) (vals []string, err error) {
	return redis.Strings(rc.Do("ZRANGE", key, startIndex, count))
}

// ZRANGE命令
func (rc *RedisClient) ZRangeWithScores(key string, startIndex, count int64) (vals map[string]int, err error) {
	return redis.IntMap(rc.Do("ZRANGE", key, startIndex, count, "WITHSCORES"))
}

// ZCARD命令
func (rc *RedisClient) ZCard(key string) (count int64, err error) {
	return redis.Int64(rc.Do("ZCARD", key))
}

// ZREVRANGEBYSCORE命令
func (rc *RedisClient) ZRevRangeByScore(key string, start int64, end int64) (vals []string, err error) {
	return redis.Strings(rc.Do("ZREVRANGEBYSCORE", key, start, end))
}

// ZADD命令
func (rc *RedisClient) ZAdd(key string, score int64, member string) (ret bool, err error) {
	return redis.Bool(rc.Do("ZADD", key, score, member))
}

// ZAdd命令
func (rc *RedisClient) ZAddList(key string, valMap map[int64]string) (err error) {
	size := 2 * len(valMap)
	arr := make([]interface{}, 0, size+1)
	arr = append(arr, key)
	for member, score := range valMap {
		arr = append(arr, score)
		arr = append(arr, member)
	}
	_, err = rc.Do("ZAdd", arr...)
	return
}

// HDEL命令
func (rc *RedisClient) HDel(key string, fields ...string) (ret int64, err error) {
	fieldsStr := make([]interface{}, 0, len(fields)+1)
	fieldsStr = append(fieldsStr, key)
	for _, field := range fields {
		fieldsStr = append(fieldsStr, field)
	}
	return redis.Int64(rc.Do("HDEL", fieldsStr...))
}

// HEXISTS命令
func (rc *RedisClient) HExists(key, field string) (ret int64, err error) {
	return redis.Int64(rc.Do("HEXISTS", key, field))
}

// HKEYS命令
func (rc *RedisClient) HKeys(key string) (ret []string, err error) {
	return redis.Strings(rc.Do("HKEYS", key))
}

// HMSET命令
func (rc *RedisClient) HMSet(key string, vals ...string) (ret string, err error) {
	if len(vals)%2 != 0 {
		return "", errors.New("field value is not coupled")
	}
	fvs := make([]interface{}, 0, len(vals)*2)
	fvs = append(fvs, key)
	for _, fv := range vals {
		fvs = append(fvs, fv)
	}
	return redis.String(rc.Do("HMSET", fvs...))
}

// HINCRBY命令
func (rc *RedisClient) HIncrBy(key, field string, val int64) (ret int64, err error) {
	return redis.Int64(rc.Do("HINCRBY", key, field, val))
}

// HGET命令
func (rc *RedisClient) HGet(key, field string) (val string, err error) {
	return redis.String(rc.Do("HGET", key, field))
}

// HMGET命令
func (rc *RedisClient) HMGet(key string, fields []string) (vals []string, err error) {
	arr := make([]interface{}, len(fields)+1)
	arr[0] = key
	for k := range fields {
		arr[k+1] = fields[k]
	}
	return redis.Strings(rc.Do("HMGET", arr...))
}

// HGETALL命令
func (rc *RedisClient) HGetAll(key string) (val map[string]string, err error) {
	return redis.StringMap(rc.Do("HGETALL", key))
}

// HSET命令
func (rc *RedisClient) HSet(key, field, val string) (ret int, err error) {
	return redis.Int(rc.Do("HSET", key, field, val))
}

// HGETALL命令
func (rc *RedisClient) MultiHGetAll(keys []string) (ret []map[string]string, err error) {
	if len(keys) == 0 {
		return
	}

	size := len(keys)
	cmds := make([]string, size)
	args := make([][]interface{}, size)
	for i := 0; i < size; i++ {
		cmds[i] = "HGETALL"
		args[i] = []interface{}{keys[i]}
	}
	var rawVals []interface{}
	var pipelineErr []error
	rawVals, pipelineErr, err = rc.Pipeline(cmds, args)
	if err != nil {
		return
	}

	ret = make([]map[string]string, 0, len(rawVals))
	for i := 0; i < len(rawVals); i++ {
		tmpVals, e := redis.StringMap(rawVals[i], pipelineErr[i])
		if e != nil {
			err = fmt.Errorf("pipeline HGETALL commands error: %v", e)
		}
		ret = append(ret, tmpVals)
	}
	return
}

// 执行命令
func (rc *RedisClient) Do(cmd string, args ...interface{}) (val interface{}, err error) {
	conn := rc.RedisPool.Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

// pipeline
func (rc *RedisClient) Pipeline(cmds []string, args [][]interface{}) (vals []interface{}, pipelineErr []error, err error) {
	if len(cmds) != len(args) {
		err = errors.New("pipeline 中的cmds 和 args 的长度不相等")
		return
	}

	cmdsLen := len(cmds)
	conn := rc.RedisPool.Get()
	defer conn.Close()

	//pipeline send
	for i := 0; i < cmdsLen; i++ {
		err = conn.Send(cmds[i], args[i]...)
		if err != nil {
			return
		}
	}

	//pipeline flush
	err = conn.Flush()
	if err != nil {
		return
	}

	//pipeline receive
	vals = make([]interface{}, 0, cmdsLen)
	pipelineErr = make([]error, 0, cmdsLen)
	for i := 0; i < cmdsLen; i++ {
		v, e := conn.Receive()
		vals = append(vals, v)
		pipelineErr = append(pipelineErr, e)
	}
	return
}

// 关闭连接池
func (rc *RedisClient) Close() {
	rc.RedisPool.Close()
}
