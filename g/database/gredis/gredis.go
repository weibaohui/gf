// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Redis客户端.
// Redis中文手册文档请参考：http://redisdoc.com/ ，Redis官方命令请参考：https://redis.io/commands
package gredis

import (
    "time"
    "github.com/gomodule/redigo/redis"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/util/gconv"
)

const (
    gDEFAULT_POOL_MAX_IDLE      = 1
    gDEFAULT_POOL_MAX_ACTIVE    = 10
    gDEFAULT_POOL_IDLE_TIMEOUT  = 180 * time.Second
    gDEFAULT_POOL_MAX_LIFE_TIME = 60  * time.Second
)

// Redis客户端
type Redis struct {
    conn redis.Conn
    pool *redis.Pool
}

// Redis链接池统计信息
type PoolStats struct {
    redis.PoolStats
}

// 连接池map
var pools = gmap.NewStringInterfaceMap()

// 创建redis操作对象
// address参数格式 host:port
func New(address string, db ... interface{}) *Redis {
    r := &Redis{}
    dialDb := 0
    if len(db) > 0 {
        dialDb = gconv.Int(db[0])
    }
    poolKey := address + "," + gconv.String(dialDb)
    if v := pools.Get(poolKey); v == nil {
        pool := &redis.Pool {
            MaxIdle         : gDEFAULT_POOL_MAX_IDLE,
            MaxActive       : gDEFAULT_POOL_MAX_ACTIVE,
            IdleTimeout     : gDEFAULT_POOL_IDLE_TIMEOUT,
            MaxConnLifetime : gDEFAULT_POOL_MAX_LIFE_TIME,
            Dial            : func() (redis.Conn, error) {
                c, err := redis.Dial("tcp", address)
                if err != nil {
                    return nil, err
                }
                c.Do("SELECT", dialDb)
                return c, nil
            },
        }
        pools.Set(poolKey, pool)
        r.pool = pool
    } else {
        r.pool = v.(*redis.Pool)
    }
    r.conn = r.pool.Get()
    return r
}

// 关闭链接
func (r *Redis) Close() error {
    return r.conn.Close()
}

// 设置属性 - MaxIdle
func (r *Redis) SetMaxIdle(value int) {
    r.pool.MaxIdle = value
}

// 设置属性 - MaxActive
func (r *Redis) SetMaxActive(value int) {
    r.pool.MaxActive = value
}

// 设置属性 - IdleTimeout
func (r *Redis) SetIdleTimeout(value time.Duration) {
    r.pool.IdleTimeout = value
}

// 设置属性 - MaxConnLifetime
func (r *Redis) SetMaxConnLifetime(value time.Duration) {
    r.pool.MaxConnLifetime = value
}

// 获取当前连接池统计信息
func (r *Redis) Stats() *PoolStats {
    return &PoolStats{r.pool.Stats()}
}

// 执行命令 - Do
func (r *Redis) Do(command string, args ...interface{}) (interface{}, error) {
    return r.conn.Do(command, args...)
}

// 执行命令 - Send
func (r *Redis) Send(command string, args ...interface{}) error {
    return r.conn.Send(command, args...)
}

