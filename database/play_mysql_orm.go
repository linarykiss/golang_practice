package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var MysqlDB *gorm.DB

type User2 struct {
	Id   int    `gorm:"size:11;primary_key;AUTO_INCREMENT;not null" json:"id"`
	Age  int    `gorm:"size:11;DEFAULT NULL" json:"age"`
	Name string `gorm:"size:255;DEFAULT NULL" json:"name"`
	//gorm后添加约束，json后为对应mysql里的字段
}

func PlayOrm() {
	MysqlDB, err := gorm.Open("mysql", "root:123456789@tcp(127.0.0.1:3306)/atest?charset=utf8")
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	} else {
		fmt.Println("connect database success")
		MysqlDB.SingularTable(true)
		MysqlDB.AutoMigrate(&User2{}) //自动建表
		fmt.Println("create table success")
	}
	defer MysqlDB.Close()

	Router()
}

func Router() {
	router := gin.Default()
	//路径映射
	router.GET("/user", InitPage)
	router.POST("/user/create", CreateUser)
	router.GET("/user/list", ListUser)
	router.PUT("/user/update", UpdateUser)
	router.GET("/user/find", GetUser)
	router.DELETE("/user/:id", DeleteUser)
	router.Run(":8080")
}

//每个路由都对应一个具体的函数操作,从而实现了对user的增,删,改,查操作
func InitPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func CreateUser(c *gin.Context) {
	var user User2
	c.BindJSON(&user)            //使用bindJSON填充对象
	MysqlDB.Create(&user)        //创建对象
	c.JSON(http.StatusOK, &user) //返回页面
}

func UpdateUser(c *gin.Context) {
	var user User2
	id := c.PostForm("id")                //post方法取相应字段
	err := MysqlDB.First(&user, id).Error //数据库查找主键=ID的第一行
	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err.Error())
	} else {
		c.BindJSON(&user)
		MysqlDB.Save(&user) //提交更改
		c.JSON(http.StatusOK, &user)
	}
}
func ListUser(c *gin.Context) {
	var user []User2
	line := c.Query("line")
	MysqlDB.Limit(line).Find(&user) //限制查找前line行
	c.JSON(http.StatusOK, &user)
}
func GetUser(c *gin.Context) {
	id := c.Query("id")
	var user User2
	err := MysqlDB.First(&user, id).Error
	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err.Error())
	} else {
		c.JSON(http.StatusOK, &user)
	}
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user User2
	MysqlDB.Where("id = ?", id).Delete(&user)
	c.JSON(http.StatusOK, gin.H{
		"data": "this has been deleted!",
	})
}
