package cache

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/astaxie/beego/cache"
)

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey = "MY_CACHE"
)

type IMyCache interface {
	GetAll() []interface{}
	GetVals() []interface{}
	GetKeys() []interface{}

	cache.Cache
}

// MyCache is Redis cache adapter.
type MyCache struct {
	p        *redis.Pool // redis connection pool
	connInfo string
	dbNum    int
	key      string
	password string
	IMyCache
}

// NewMyCache create new redis cache with default collection name.
func NewMyCache() cache.Cache {
	return &MyCache{key: DefaultKey}
}

// actually do the redis cmds
func (rc *MyCache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

func (rc *MyCache) GetVals() []interface{} {
	if v, err := rc.do("HVALS", rc.key); err == nil {
		if all, done := v.([]interface{}); done {
			return all
		}
	}
	return nil
}

func (rc *MyCache) GetKeys() []interface{} {
	if v, err := rc.do("HKEYS", rc.key); err == nil {
		if all, done := v.([]interface{}); done {
			return all
		}
	}
	return nil
}

// Get all cache from redis.
func (rc *MyCache) GetAll() []interface{} {
	if v, err := rc.do("HGETALL", rc.key); err == nil {
		if all, done := v.([]interface{}); done {
			return all
		}
	}
	return nil
}

// Get cache from redis.
func (rc *MyCache) Get(key string) interface{} {
	if v, err := rc.do("HGET", rc.key, key); err == nil {
		return v
	}
	return nil
}

// GetMulti get cache from redis.
func (rc *MyCache) GetMulti(keys []string) []interface{} {
	args := make([]interface{}, len(keys)+1)
	args[0] = rc.key

	for i := range args {
		args[i+1] = keys[i]
	}

	if v, err := rc.do("HMGET", args); err == nil {
		if res, ok := v.([]interface{}); ok {
			return res
		}
	}
	return nil
}

// Put put cache to redis.
func (rc *MyCache) Put(key string, val interface{}, timeout time.Duration) error {
	var err error

	if _, err = rc.do("HSET", rc.key, key, val); err != nil {
		return err
	}

	timer := time.NewTimer(timeout)
	go func() {
		<-timer.C
		rc.Delete(key)
	}()

	return err
}

// Delete delete cache in redis.
func (rc *MyCache) Delete(key string) error {
	var err error
	_, err = rc.do("HDEL", rc.key, key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *MyCache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("HEXISTS", rc.key, key))
	if err != nil {
		return v
	}
	return false
}

// Incr increase counter in redis.
func (rc *MyCache) Incr(key string) error {
	_, err := redis.Bool(rc.do("HINCRBY", rc.key, key, 1))
	return err
}

// Decr decrease counter in redis.
func (rc *MyCache) Decr(key string) error {
	_, err := redis.Bool(rc.do("HINCRBY", rc.key, key, -1))
	return err
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *MyCache) ClearAll() error {
	cachedKeys, err := redis.Strings(rc.do("HKEYS", rc.key))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = rc.do("HDEL", rc.key, str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *MyCache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	rc.key = cf["key"]
	rc.connInfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *MyCache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.connInfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selectErr := c.Do("SELECT", rc.dbNum)
		if selectErr != nil {
			c.Close()
			return nil, selectErr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func init() {
	cache.Register("my", NewMyCache)
}
