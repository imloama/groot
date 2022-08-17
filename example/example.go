package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imloama/groot"
)

func main() {
	server := groot.NewDefaultServer()
	r := server.Engine
	r.GET("/", func(ctx *gin.Context) {
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

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
