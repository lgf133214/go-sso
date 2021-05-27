package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-sso/internal/conf"
	"log"
	"strconv"
	"time"
)

var db *redis.Client
var prefix = conf.Cfg.Redis.Prefix

func init() {
	db = redis.NewClient(&redis.Options{
		Addr:     conf.Cfg.Redis.Addr + ":" + strconv.FormatInt(conf.Cfg.Redis.Port, 10),
		Username: conf.Cfg.Redis.UserName,
		Password: conf.Cfg.Redis.Password,
		DB:       conf.Cfg.Redis.DB,
	})
	err := db.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func AppIDIsExist(appID string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2e8)
	defer cancel()

	key := prefix + ":app-id"
	return db.HExists(ctx, key, appID).Val()
}

func AddAppIDs(appIDs []interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	err := db.Del(ctx, prefix+":app-id").Err()
	if err != nil {
		log.Fatal(err)
	}

	err = db.HMSet(ctx, prefix+":app-id", appIDs...).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func GetAppRedirect(appId string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2e8)
	defer cancel()

	result, err := db.HGet(ctx, prefix+":app-id", appId).Result()
	if err != nil {
		return "", false
	}
	return result, true
}

func GetUuidByAT(at string) (uuid string, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2e8)
	defer cancel()

	uuid = db.Get(ctx, prefix+":access-token:"+at).String()
	if uuid != "" {
		ok = true
	}
	return
}

func GetATExpire(at string) (expire time.Duration, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2e8)
	defer cancel()

	expire = db.TTL(ctx, prefix+":access-token:"+at).Val()
	if expire != 0 {
		ok = true
	}
	return
}

func SetAT(at, uuid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2e8)
	defer cancel()

	err := db.SetEX(ctx, prefix+":access-token:"+at, uuid, conf.Cfg.Redis.Expire).Err()
	if err != nil {
		log.Println("imposable case occur")
		log.Println(err)
	}
}

func GetUuidByST(st string) (uuid string, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2e8)
	defer cancel()

	uuid = db.Get(ctx, prefix+":service-ticket:"+st).String()
	if uuid != "" {
		ok = true
	}
	return
}

func GetSTExpire(st string) (expire time.Duration, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2e8)
	defer cancel()

	expire = db.TTL(ctx, prefix+":service-ticket:"+st).Val()
	if expire != 0 {
		ok = true
	}
	return
}

func SetSTExpire(st, uuid string, ex time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 2e8)
	defer cancel()

	err := db.SetEX(ctx, prefix+":service-ticket:"+st, uuid, ex).Err()
	if err != nil {
		log.Println("imposable case occur")
		log.Println(err)
	}
}
