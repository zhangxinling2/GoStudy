package main

import (
	"GoStudy/dataStore/fatRank"
	"GoStudy/httpServer/gin/practice/moments"
	"GoStudy/httpServer/httpPratice/frinterface"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

//type PersonalMomentData struct {
//	Id int64 `json:"id,omitempty" gorm:"primaryKey;column:id"`
//
//	PersonId int64 `json:"personId,omitempty" gorm:"column:person_id"`
//
//	CreatedTime time.Time `json:"createdTime,omitempty" gorm:"column:created_time"`
//
//	Content string `json:"content,omitempty" gorm:"column:content"`
//
//	Fatrate float32 `json:"fatrate,omitempty" gorm:"column:fatrate"`
//
//	Visible bool `json:"visible,omitempty" gorm:"column:visible"`
//}

func ConnectDataBase() (*gorm.DB, error) {
	conn, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/mysql?parseTime=true"))
	if err != nil {
		log.Fatal("数据库连接失败")
	}
	return conn, nil
}
func main() {
	conn, err := ConnectDataBase()
	if err != nil {
		log.Println("连接数据库失败")
		return
	}
	person := fatRank.PersonalInformation{}

	conn.First(&person, 1)
	if err = conn.AutoMigrate(&fatRank.PersonalMoment{}); err != nil {
		log.Println("建表失败", err)
		return
	}
	var db frinterface.ServeInterface = NewDbRank(conn, NewFatRateRank())
	//db := NewDbRank(conn, NewFatRateRank())
	var moment moments.Moments = NewMomentServer(conn, person)
	if initRank, ok := db.(frinterface.RankInitInterface); ok {
		if err := initRank.Init(); err != nil {
			log.Fatal("初始化失败", err)
		}
	}
	r := gin.Default()
	pprof.Register(r)

	r.POST("/registry", func(c *gin.Context) {
		pi := fatRank.PersonalInformation{}
		if err := c.BindJSON(&pi); err != nil {
			c.JSON(400, gin.H{
				"ErrMessage": "无法解析注册信息",
			})
		}
		if err := db.RegisterPersonInformation(&pi); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ErrMessage": fmt.Sprintf("无法注册信息,%v", err),
			})
			return
		}
		c.JSON(200, gin.H{
			"ErrMessage": "success",
		})
	})
	r.PUT("/personinfo", func(c *gin.Context) {
		pi := fatRank.PersonalInformation{}
		if err := c.BindJSON(&pi); err != nil {
			c.JSON(400, gin.H{
				"ErrMessage": "无法解析注册信息",
			})
		}
		if fr, err := db.UpdatePersonInformation(&pi); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ErrMessage": fmt.Sprintf("无法更改信息,%v", err),
			})
			return
		} else {
			c.JSON(http.StatusOK, fr)
		}
	})
	r.GET("/rank/:name", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"ErrMessage": "name未设置",
			})
			return
		}
		if fr, err := db.GetFatrate(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ErrMessage": fmt.Sprintf("无法获取排行数据,%v", err),
			})
			return
		} else {
			c.JSON(200, fr)
		}
	})
	r.GET("/rankTop", func(c *gin.Context) {
		if frTop, err := db.GetTop(); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("无法获取排行数据,%v", err))
			return
		} else {
			c.JSON(http.StatusOK, frTop)
		}
	})
	r.POST("/moment/release", func(c *gin.Context) {
		var text string
		c.BindJSON(&text)
		if res, err := moment.ReleaseMoment(text); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("无法发布,%v", err))
			return
		} else {
			c.JSON(200, res)
		}
	})
	r.PUT("/moment/del", func(c *gin.Context) {
		var id int
		c.BindJSON(&id)
		if err = moment.DeleteMoment(int64(id)); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("无法删除,%v", err))
			return
		}
		c.JSON(200, "删除成功")
	})
	r.GET("/moment", func(c *gin.Context) {
		mo, err := moment.GetAllMoment()
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("无法得到,%v", err))
			return
		}
		c.JSON(200, mo)
	})
	r.Run(":8080")
	// http.ListenAndServe(":8088", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// time.Sleep(2 * time.Second)
	// 	// w.Write([]byte("hello 你好"))

	// 	//request读的时候要使用POST
	// 	if r.Body == nil {
	// 		w.Write([]byte("no body"))
	// 		return
	// 	}
	// 	data, _ := ioutil.ReadAll(r.Body)
	// 	defer r.Body.Close()
	// 	encoded := base64.StdEncoding.EncodeToString(data)
	// 	w.Write(append(data, []byte(encoded)...))

	// 	// qp := r.URL.Query()
	// 	// data, _ = json.Marshal(qp)
	// 	// w.Write([]byte("hello 你好" + string(data)))
	// }))
}
