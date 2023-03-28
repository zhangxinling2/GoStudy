package main

import (
	"GoStudy/dataStore/fatRank"
	"encoding/base64"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	pprof.Register(r)
	r.GET("/history", func(c *gin.Context) {
		name := c.Query("name")
		c.JSON(http.StatusOK, []*fatRank.PersonalInformation{
			{
				Id:     0,
				Name:   name,
				Sex:    "男",
				Tall:   0,
				Weight: 0,
				Age:    0,
			},
		})
	})
	r.GET("/history/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{
			"spacial": base64.StdEncoding.EncodeToString([]byte(name)),
		})
	})
	r.POST("/register", func(c *gin.Context) {
		pi := &fatRank.PersonalInformation{}
		if err := c.BindJSON(pi); err != nil {
			c.JSON(400, gin.H{
				"message": "无法读取注册信息",
			})
			return
		}
		//todo :注册到排行榜
		c.JSON(200, "")
	})
	r.Run(":8080")
}
