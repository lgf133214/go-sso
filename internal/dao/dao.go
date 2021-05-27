package dao

import (
	"go-sso/internal/dao/mysql"
	"go-sso/internal/dao/redis"
	"go-sso/pkg/utils/uuid"
	"time"
)

func init() {
	redis.AddAppIDs(mysql.AppIDs())
}

func GenAT(email string) string {
	s := mysql.GetUuidByEmail(email)
	at := uuid.GenUuid()
	redis.SetAT(at, s)
	return at
}

func GetRedirectUrl(appID string) (url string, ok bool) {
	return redis.GetAppRedirect(appID)
}

func GenST(at string) (string, time.Duration, bool) {
	u, ok := redis.GetUuidByAT(at)
	if !ok {
		return "", 0, false
	}
	st := uuid.GenUuid()
	ex, _ := redis.GetATExpire(at)
	redis.SetSTExpire(st, u, ex)
	return st, ex, true
}

func VerifyUser(mail, password, ip string) bool {
	_, ok := mysql.UserIsValid(mail, password, ip)
	return ok
}
