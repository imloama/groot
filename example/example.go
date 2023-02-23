package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imloama/groot"
)

func main() {
	server := groot.NewDefaultServer()
	r := server.Engine
	r.Use(groot.TraceIdMiddleware())
	r.GET("/", func(ctx *gin.Context) {
		groot.TDebug(ctx, "hello,world!==============")
		ctx.String(http.StatusOK, "hello, world!")
	})

	r.GET("/user", func(ctx *gin.Context) {
		db := groot.GetDb()
		if db == nil {
			ctx.String(http.StatusOK, "no orm!")
			return
		}
		var mapdata = make(map[string]interface{})
		db.Raw(`select * from sys_user where user_id = 1`).Scan(&mapdata)
		ctx.JSON(http.StatusOK, mapdata)
	})

	r.GET("/redis", func(ctx *gin.Context) {
		db := groot.GetRedisClient()
		db.SetEx(ctx.Request.Context(), "a", "hello, redis!", time.Second*20)
		result := db.Get(ctx.Request.Context(), "a").Val()
		ctx.JSON(http.StatusOK, result)
	})

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
