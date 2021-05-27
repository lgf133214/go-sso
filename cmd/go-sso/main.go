package main

import (
	"github.com/gin-gonic/gin"
	"go-sso/handler"
	"go-sso/internal/conf"
	"go-sso/middleware"
	"log"
	"net/http"
)

func init() {
	if conf.Cfg.Gin.Release {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	r := gin.Default()
	{
		r.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})
	}

	r.LoadHTMLFiles("template/login.html", "template/register.html")
	r.StaticFS("/static", http.Dir("template/static"))

	app := r.Group("/app-id/:app-id")
	{
		app.Use(middleware.ValidateAppID)

		app.GET("/login", handler.LoginGET)
		app.POST("/login", handler.LoginPOST)

		app.POST("/logout", handler.Logout)

		app.GET("/validate/:st", handler.ValidateST)
	}
	{
		r.GET("/register", handler.RegisterGET)
		r.POST("/register", handler.RegisterPOST)
		r.GET("/activate", handler.Activate)
	}

	log.Fatal(r.Run(conf.GetListenAddr()))
}
