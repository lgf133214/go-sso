package handler

import (
	"github.com/gin-gonic/gin"
	"go-sso/internal/conf"
	"go-sso/internal/dao"
	"go-sso/internal/dao/model"
	"go-sso/internal/dao/mysql"
	"go-sso/internal/dao/redis"
	"go-sso/pkg/utils/hash"
	"go-sso/pkg/utils/mail"
	"go-sso/pkg/utils/random"
	"go-sso/pkg/utils/uuid"
	"net/url"
	"strconv"
	"time"
)

func LoginGET(c *gin.Context) {
	redirect := c.Query("redirect")
	if redirect == "" {
		c.AbortWithStatus(400)
		return
	}

	at, _ := c.Cookie("Access-Token")
	if at != "" {
		if ex, ok := redis.GetATExpire(at); ok {
			if ex > time.Minute*5 {
				st, _, _ := dao.GenST(at)

				s, _ := dao.GetRedirectUrl(c.Param("app-id"))
				uri, _ := url.Parse(s)
				q := uri.Query()
				q.Add("st", st)
				q.Add("redirect", redirect)
				uri.RawQuery = q.Encode()

				c.Redirect(302, uri.String())
				return
			}
		}
		c.SetCookie("Access-Token", "", 0, "/", conf.Cfg.Gin.Host, true, true)
	}

	c.HTML(302, "login.html", gin.H{
		"Redirect": redirect,
	})
}

func LoginPOST(c *gin.Context) {
	redirect := c.PostForm("redirect")
	if redirect == "" {
		c.AbortWithStatus(400)
		return
	}

	email := c.PostForm("email")
	if email == "" {
		c.AbortWithStatus(400)
		return
	}

	password := c.PostForm("password")
	if password == "" {
		c.AbortWithStatus(400)
		return
	}

	if !dao.VerifyUser(email, password, c.ClientIP()) {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "用户名或密码错误",
		})
		return
	}

	at := dao.GenAT(email)
	c.SetCookie("Access-Token", at, int(conf.Cfg.Redis.Expire.Seconds()), "/", conf.Cfg.Gin.Host, true, true)

	st, ex, _ := dao.GenST(at)
	uri, _ := dao.GetRedirectUrl(c.Param("app-id"))
	uri += "?st=" + st
	uri += "&ex=" + strconv.FormatInt(int64(ex), 10)
	uri += "&redirect=" + redirect
	c.Redirect(302, "http://"+uri)
}

func Logout(c *gin.Context) {

}

func ValidateST(c *gin.Context) {
	st := c.Param("st")
	if st == "" {
		c.AbortWithStatus(400)
		return
	}

	uid, ok := redis.GetUuidByST(st)
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	ex, ok := redis.GetSTExpire(st)
	if !ok {
		c.AbortWithStatus(400)
		return
	}
	c.JSON(200, gin.H{
		"expire": ex,
		"uuid":   uid,
	})
}

func RegisterGET(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}

func RegisterPOST(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(400, gin.H{
			"msg": "未填写邮箱",
		})
		return
	}
	if !mail.IsValid(email) {
		c.JSON(400, gin.H{
			"msg": "邮箱格式错误",
		})
		return
	}

	password := c.PostForm("password")
	if password == "" {
		c.JSON(400, gin.H{
			"msg": "未填写密码",
		})
		return
	}

	if mysql.IsExistByEmail(email) {
		c.JSON(400, gin.H{
			"msg": "该邮箱已被注册",
		})
		return
	}
	salt := random.GetRandomString(4)
	code := random.GetRandomString(20)
	user := model.User{Uuid: uuid.GenUuid(), Email: email, Password: hash.MD5(password + salt), Salt: salt, IP: c.ClientIP(), Status: 0, VerifyCode: code}
	mysql.Register(user)
	mail.SendMail(email, "点击激活链接激活",
		"http://"+conf.Cfg.Gin.Host+":"+strconv.FormatInt(conf.Cfg.Gin.Port, 10)+
			"/activate?uuid="+user.Uuid+"&verify_code="+code)
	c.JSON(200, gin.H{
		"msg": "验证邮件已发送",
	})
}

func Activate(c *gin.Context) {
	id := c.Query("uuid")
	code := c.Query("verify_code")
	if id == "" || code == "" {
		c.AbortWithStatus(400)
		return
	}
	if mysql.IsActivate(id) {
		c.AbortWithStatusJSON(200, gin.H{
			"msg": "已激活",
		})
		return
	}

	verifyCode, ok := mysql.VerifyCodeByUuid(id)
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	if verifyCode != code {
		c.AbortWithStatus(400)
		return
	}
	mysql.Activate(id, c.ClientIP())
	c.JSON(200, gin.H{
		"msg": "激活成功",
	})
}
