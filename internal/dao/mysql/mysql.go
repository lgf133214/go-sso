package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-sso/internal/conf"
	"go-sso/internal/dao/model"
	"go-sso/pkg/utils/hash"
	"log"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conf.Cfg.Mysql.UserName, conf.Cfg.Mysql.Password,
		conf.Cfg.Mysql.Addr, conf.Cfg.Mysql.Port, conf.Cfg.Mysql.DB))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func AppIDIsValid(appID string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	err := db.QueryRowContext(ctx, "select `app_id` from `apps` where `app_id`=? and `status`=1", appID).Scan(&appID)

	// ignore error cause we don't care it :)
	if err != nil {
		return false
	}
	return true
}

func AppIDs() (ret []interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 1e9)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select `app_id`, `redirect` from `apps` where `status`=1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tmp1, tmp2 string
		err := rows.Scan(&tmp1, &tmp2)
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, tmp1, tmp2)
	}
	return
}

func UserIsValid(email, password, ip string) (string, bool) {
	salt, ok := GetSalt(email)
	if !ok {
		return "", false
	}
	password = hash.MD5(password + salt)
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	var uuid string
	err := db.QueryRowContext(ctx, "select `uuid` from `users` where `email`=? and `password`=? and `status`=1", email, password).Scan(&uuid)
	fmt.Println("select `uuid` from `users` where `email`=? and `password`=? and `stasus`=1", email, password)
	if err != nil || uuid == "" {
		return "", false
	}
	go traceRecord(model.Trace{Uuid: uuid, Type: 2, Ip: ip})
	return uuid, true
}

func GetSalt(email string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	var s string
	if db.QueryRowContext(ctx, "select `salt` from `users` where `email`=?", email).Scan(&s) != nil {
		return "", false
	}
	return s, true
}

func IsExistByEmail(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	err := db.QueryRowContext(ctx, "select `email` from `users` where `email`=?", email).Scan(&email)
	// ignore error
	if err != nil {
		return false
	}
	return true
}

func VerifyCodeByUuid(uuid string) (verifyCode string, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	err := db.QueryRowContext(ctx, "select `verify_code` from `users` where `uuid`=?", uuid).Scan(&verifyCode)
	if err != nil {
		log.Println("imposable case occur")
		log.Println(err)
		return
	}
	return verifyCode, true
}

func Activate(uuid, ip string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	_, err := db.ExecContext(ctx, "update `users` set `status`=1 where `uuid`=?", uuid)
	if err != nil {
		log.Println("imposable case occur")
		log.Println(err)
		return false
	}
	go traceRecord(model.Trace{Uuid: uuid, Type: 1, Ip: ip})
	return true
}

func Register(user model.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	_, err := db.ExecContext(ctx, "insert into `users`(`uuid`, `email`, `password`, `salt`, `status`, `verify_code`) values (?,?,?,?,?,?)",
		user.Uuid, user.Email, user.Password, user.Salt, user.Status, user.VerifyCode)
	if err != nil {
		log.Println("imposable case occur")
		log.Println(err)
		return false
	}

	go traceRecord(model.Trace{Uuid: user.Uuid, Type: 0, Ip: user.IP})
	return true
}

func traceRecord(trace model.Trace) {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	_, err := db.ExecContext(ctx, "insert into `traces`(`uuid`, `type`, `ip`) values (?,?,?)",
		trace.Uuid, trace.Type, trace.Ip)
	if err != nil {
		log.Println("imposable case occur")
		log.Println(err)
	}
}

func GetUuidByEmail(email string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	err := db.QueryRowContext(ctx, "select `uuid` from `users` where `email`=?", email).Scan(&email)
	if err != nil {
		return ""
	}
	return email
}

func IsActivate(uid string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5e8)
	defer cancel()

	var status int
	err := db.QueryRowContext(ctx, "select `status` from `users` where `uuid`=?", uid).Scan(&status)
	if err != nil {
		return false
	}
	return status == 1
}
