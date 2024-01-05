// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package redis_client

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/leovs/yc-go-sdk/log"
	"reflect"
	"time"
)

type RedisClient struct {
	Conn           redis.Conn
	Pool           *redis.Pool
	masterName     string // Sentinel 哨兵模式 Master名字
	address        string // 地址 localhost:6379
	password       string // 密码
	dbIds          int    // redisDB
	maxIdle        int    // redis连接池最大空闲连接数
	maxActive      int    // redis连接池最大激活连接数, 0为不限制
	connectTimeout int    // redis连接超时时间, 单位毫秒
	readTimeout    int    // redis读取超时时间, 单位毫秒
	writeTimeout   int    // redis写入超时时间, 单位毫秒
	masterAddr     string // master地址
}

func (e *RedisClient) InitRedis(
	masterName string, address string, password string, dbIds int, maxIdle int,
	maxActive int, connectTimeout int, readTimeout int, writeTimeout int,
) {
	e.dbIds = dbIds
	e.masterName = masterName
	e.address = address
	e.password = password
	e.maxIdle = maxIdle
	e.maxActive = maxActive
	e.connectTimeout = connectTimeout
	e.readTimeout = readTimeout
	e.writeTimeout = writeTimeout
	e.initRedisPool()
}

// 初始化哨兵模式
func (e *RedisClient) initSentinel(masterName string) (err error) {
	e.Conn, err = redis.Dial("tcp", e.address,
		redis.DialConnectTimeout(time.Duration(e.connectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(e.readTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(e.writeTimeout)*time.Millisecond),
	)
	if err != nil {
		log.Error("#50000@初始化Redis失败 %v", err.Error())
		return err
	}
	defer e.Conn.Close()
	res, err := redis.Strings(e.Conn.Do("sentinel", "get-master-addr-by-name", masterName))
	if err != nil {
		log.Error("#50001@获取master信息失败 %v", err.Error())
		return err
	}

	if len(res) < 2 {
		log.Error("#50002@master地址信息错误")
		return errors.New("master地址信息错误")
	}
	e.masterAddr = res[0] + ":" + res[1]
	return nil
}

// 初始化连接池
func (e *RedisClient) initRedisPool() {
	e.Pool = &redis.Pool{
		MaxIdle:      e.maxIdle,
		MaxActive:    e.maxActive,
		IdleTimeout:  300 * time.Second,
		Dial:         e.redisDial,
		TestOnBorrow: e.redisTestOnBorrow,
		Wait:         true,
	}
}

// 连接redis
func (e *RedisClient) redisDial() (redis.Conn, error) {
	// 初始化哨兵信息
	err1 := e.initSentinel(e.masterName)
	if err1 != nil {
		log.Error("#50003@初始化哨兵失败 %v", err1.Error())
		return nil, err1
	}

	conn, err := redis.Dial(
		"tcp",
		e.masterAddr,
		redis.DialConnectTimeout(time.Duration(e.connectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(e.readTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(e.writeTimeout)*time.Millisecond),
	)
	if err != nil {
		log.Error("#50004@连接redis失败 %v", err.Error())
		return nil, err
	}

	if e.password != "" {
		if _, err := conn.Do("AUTH", e.password); err != nil {
			_ = conn.Close()
			log.Error("#50005@redis认证失败 %v", err.Error())
			return nil, err
		}
	}

	//获取master信息
	_, err = conn.Do("SELECT", e.dbIds)
	if err != nil {
		_ = conn.Close()
		log.Error("#50006@redis初始化失败 %v", err.Error())
		return nil, err
	}

	return conn, nil
}

// 从池中取出连接后，判断连接是否有效
func (e *RedisClient) redisTestOnBorrow(conn redis.Conn, _ time.Time) error {
	_, err := conn.Do("PING")
	if err != nil {
		log.Error("#50007@从redis连接池取出的连接无效 %v", err.Error())
	}
	return err
}

// Exec 执行redis命令, 执行完成后连接自动放回连接池
func (e *RedisClient) Exec(command string, args ...interface{}) (any, error) {
	if e.Pool != nil {
		conn := e.Pool.Get()
		if conn != nil {
			defer conn.Close()
			return conn.Do(command, args...)
		}
	}
	log.Error("#50007@redis 不可用 %v\n")
	return nil, errors.New("redis不可用")
}

func (e *RedisClient) Set(key string, val any) error {
	return e.SetEx(key, val, 0)
}

func (e *RedisClient) SetEx(key string, val any, expires int64) error {
	b, err := e.serialize(val)
	if err != nil {
		log.Error("SetEx serialize err:%v\n", err)
		return err
	}
	if expires > 0 {
		_, err = e.Exec("SETEX", key, expires, b)
	} else {
		_, err = e.Exec("SET", key, b)
	}
	return err
}

func (e *RedisClient) GetObject(key string, ptr any) error {
	raw, err := e.Exec("GET", key)
	if err != nil || raw == nil {
		return errors.New("empty")
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		log.Error("#50008@GetObject failed: %s", err)
		return nil
	}
	return e.deserialize(item, ptr)
}

func (e *RedisClient) serialize(value any) ([]byte, error) {
	err := e.registerGobConcreteType(value)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(value).Kind() == reflect.Struct {
		return nil, fmt.Errorf("serialize func only take pointer of a struct")
	}

	marshal, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func (e *RedisClient) deserialize(byt []byte, ptr any) (err error) {
	return json.Unmarshal(byt, ptr)
}

func (e *RedisClient) registerGobConcreteType(value any) error {
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Ptr:
		v := reflect.ValueOf(value)
		i := v.Elem().Interface()
		gob.Register(&i)
	case reflect.Struct, reflect.Map, reflect.Slice:
		gob.Register(value)
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
	default:
		return fmt.Errorf("unhandled type: %v", t)
	}
	return nil
}
