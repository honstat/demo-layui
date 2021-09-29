package main

import (
	"net"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/pprof"
	"strings"
)

func main(){
	Run()
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
func Run() {

	router := gin.Default()
	//跨域
	router.Use(Cors())

	router.POST("/user/list", ListUser)
	router.POST("/user/remove", RemoveUser)
	router.POST("/user/add", UserAdd)
	pprof.Register(router)

	//启动参数
	_ = router.Run(net.JoinHostPort("127.0.0.1", "8000"))
}


//数据部分
var users []User

func fillUser() []User {
	arr := make([]User, 0)
	u1 := User{Name: "张三", Sex: 1, No: 16, Birthday: "2002-01-01"}
	u2 := User{Name: "李四", Sex: 1, No: 17, Birthday: "2003-01-01"}
	u3 := User{Name: "王五", Sex: 2, No: 18, Birthday: "2004-01-01"}
	u4 := User{Name: "赵云", Sex: 1, No: 19, Birthday: "2005-01-01"}
	u5 := User{Name: "李云龙", Sex: 1, No: 20, Birthday: "2006-01-01"}
	u6 := User{Name: "张三丰", Sex: 1, No: 21, Birthday: "2007-01-01"}
	u7 := User{Name: "巴马", Sex: 1, No: 22, Birthday: "2008-01-02"}
	u8 := User{Name: "德邦", Sex: 1, No: 23, Birthday: "2008-01-01"}
	u9 := User{Name: "林动", Sex: 1, No: 24, Birthday: "2018-01-01"}
	u10 := User{Name: "唐三", Sex: 1, No: 25, Birthday: "2010-01-01"}
	u11 := User{Name: "Mac超", Sex: 3, No: 26, Birthday: "2011-01-01"}
	u12 := User{Name: "周三", Sex: 2, No: 27, Birthday: "2001-01-01"}
	arr = append(arr, u1, u2, u3, u4, u5, u6, u7, u8, u9, u10, u11, u12)
	return arr
}

type User struct {
	Name     string `json:"name" form:"name"`
	Sex      int    `json:"sex" form:"sex"`
	No       int    `json:"no" form:"no"`
	Birthday string `json:"birthday" form:"birthday"`
}

func ListUser(c *gin.Context) {
	var input User
	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "input error",
		})
		return
	}
	users = fillUser()
	//数据过滤
	if len(input.Name) > 0 {
		arr := make([]User, 0)
		for _, user := range users {
			if strings.Contains(user.Name, input.Name) {
				arr = append(arr, user)
			}
		}
		users = arr
	}
	if input.Sex > 0 {
		arr := make([]User, 0)
		for _, user := range users {
			if user.Sex == input.Sex {
				arr = append(arr, user)
			}
		}
		users = arr
	}
	c.JSON(http.StatusOK, gin.H{
		"code":        0,
		"data":        users,
		"total_count": len(users),
		"msg":         "success",
	})
}
func UserAdd(c *gin.Context) {

	var input User
	err := c.ShouldBind(&input)
	if err != nil || input.No <= 0 || len(input.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "input error",
		})
		return
	}
	//添加
	users = append(users, input)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": true,
		"msg":  "success",
	})

}
func RemoveUser(c *gin.Context) {
	type Re struct {
		UserNo int `json:"no" form:"no"`
	}
	var input = Re{}
	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "input error",
		})
		return
	}
	if len(users) == 0 {
		users = fillUser()
	}
	deleteIndex := -1
	for i, user := range users {
		if user.No == input.UserNo {
			deleteIndex = i
			break
		}
	}
	if deleteIndex >= 0 {
		users = append(users[:deleteIndex], users[deleteIndex+1:]...)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": true,
			"msg":  "success",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "数据不存在",
		})
	}

}
