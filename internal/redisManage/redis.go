package redisManage

import (
	"context"
	"fmt"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	goredis "github.com/redis/go-redis/v9"
	"sobotctl/global"
	"strconv"
	"strings"
)

const (
	keyTypeString = "string"
	keyTypeList   = "list"
	keyTypeHash   = "hash"
	keyTypeSet    = "set"
)

type RedisKeyInfo struct {
	Key     string `json:"key"`
	KeyType string `json:"key_type"`
	TTL     int64  `json:"ttl"`
}

type RedisConfig struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	ReLoad bool   `json:"reload"`
	Desc   string `json:"desc"`
}

func (rc *RedisConfig) SetDesc() {
	rc.Desc = ""
}

func (rc *RedisConfig) SetReLoad() {
	rc.ReLoad = false
}

type SRedis struct {
	RedisGo redisgo.Conn
	Goredis *goredis.Client
}

func SetUpRedisGo() (redisgo.Conn, error) {
	password := redisgo.DialPassword(global.RedisSetting.Pass)
	addrStr := fmt.Sprintf("%s:%d", global.RedisSetting.IP, global.RedisSetting.Port)
	return redisgo.Dial("tcp", addrStr, password)
}

func SetUpGoRedis(db int) (*goredis.Client, error) {
	Addr := fmt.Sprintf("%s:%d", global.RedisSetting.IP, global.RedisSetting.Port)
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     Addr,
		Password: global.RedisSetting.Pass,
		DB:       db,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func SetUpGoRedisCluster(addr []string, password string) (*goredis.ClusterClient, error) {
	rdb := goredis.NewClusterClient(&goredis.ClusterOptions{
		Addrs:    addr,
		Password: password,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func NewRedisGo() (*SRedis, error) {
	conn, err := SetUpRedisGo()
	if err != nil {
		return nil, err
	}
	return &SRedis{RedisGo: conn}, nil
}

func NewGoRedis(db int) (*SRedis, error) {
	client, err := SetUpGoRedis(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &SRedis{
		Goredis: client,
	}, nil
}

func (sr *SRedis) Close() {
	if sr.Goredis != nil {
		sr.Goredis.Close()
	}
	if sr.RedisGo != nil {
		sr.RedisGo.Close()
	}
}

//func (sr *SRedis) SlowLog() ([]redisgo.SlowLog, error) {
//	result, err := sr.RedisGo.Do("SLOWLOG", "GET")
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	return redisgo.SlowLogs(result, err)
//}

func (sr *SRedis) SearchKey(keyword string) ([]*RedisKeyInfo, error) {
	respKeys := make([]*RedisKeyInfo, 0)
	keys, err := sr.Goredis.Keys(context.Background(), keyword).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var keyInfoSplit []string
	if len(keys) > 0 {
		keyInfoLua := `
			local result = {}
			-- KEYS[1]为第1个参数，lua数组下标从1开始
			local ttl = redis.call('ttl', KEYS[1]);
			local keyType = redis.call('type', KEYS[1]);
			for i = 1, #KEYS do
				local ttl = redis.call('ttl', KEYS[i]);
				local keyType = redis.call('type', KEYS[i]);
				table.insert(result, string.format("%d,%s", ttl, keyType['ok']));
			end;
			return table.concat(result, ".");`
		// 通过lua获取 ttl,type.ttl2,type2格式，以便下面切割获取ttl和type。避免多次调用ttl和type函数
		keyInfos, _ := sr.Goredis.Eval(context.Background(), keyInfoLua, keys).Result()
		keyInfoSplit = strings.Split(keyInfos.(string), ".")
	}

	for i, k := range keys {
		ttlType := strings.Split(keyInfoSplit[i], ",")
		ttl, _ := strconv.Atoi(ttlType[0])
		rk := &RedisKeyInfo{
			Key:     k,
			KeyType: ttlType[1],
			TTL:     int64(ttl),
		}
		respKeys = append(respKeys, rk)
	}
	return respKeys, nil
}

func (sr *SRedis) keyType(key string) (string, error) {
	value, err := sr.Goredis.Type(context.Background(), key).Result()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return value, nil
}

func (sr *SRedis) GetStringValue(key string) (string, error) {
	value, err := sr.Goredis.Get(context.Background(), key).Result()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return value, nil
}

func (sr *SRedis) GetHashValue(key string) (map[string]string, error) {
	value, err := sr.Goredis.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return value, nil
}

func (sr *SRedis) GetListValue(key string) ([]string, error) {
	value, err := sr.Goredis.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return value, nil
}

func (sr *SRedis) GetSetValue(key string) ([]string, error) {
	value, err := sr.Goredis.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return value, nil
}

func (sr *SRedis) GetConfig() ([]*RedisConfig, error) {
	configs := make([]*RedisConfig, 0)
	ret, err := sr.Goredis.ConfigGet(context.Background(), "*").Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for key, value := range ret {
		if key == "requirepass" {
			continue
		}
		config := &RedisConfig{
			Name:  key,
			Value: value,
		}
		configs = append(configs, config)
	}
	return configs, nil
}
