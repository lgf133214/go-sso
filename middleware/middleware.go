package middleware

import (
	"github.com/gin-gonic/gin"
	"go-sso/internal/dao/redis"
	"log"
)

func ValidateAppID(c *gin.Context) {
	appID := c.Param("app-id")
	if appID==""{
		c.AbortWithStatus(400)
		return
	}

	ok:=redis.AppIDIsExist(appID)
	if ok{
		return
	}
	c.AbortWithStatus(400)

	log.Println("app-id missed")
}
