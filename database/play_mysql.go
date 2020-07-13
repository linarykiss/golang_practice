package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // 没有导入会报错
)

/*
1、实现mysql的建表，数据的增删改查等基础操作
*/

//数据库连接信息
const (
	USERNAME = "root"
	PASSWORD = "123456789"
	NETWORK  = "tcp"
	SERVER   = "localhost"
	PORT     = 3306
	DATABASE = "atest"
)

//user表结构体定义
type User1 struct {
	Id         int    `json:"id" form:"id"`
	Username   string `json:"username" form:"username"`
	Password   string `json:"password" form:"password"`
	Status     int    `json:"status" form:"status"` // 0 正常状态， 1删除
	Createtime int64  `json:"createtime" form:"createtime"`
}

func PlayBase() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	DB, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}
	defer DB.Close()

	DB.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超时的连接就close
	DB.SetMaxOpenConns(100)                  //设置最大连接数

	//创建表
	CreateTable(DB)
	//插入数据
	InsertData(DB)
	//查询数据
	QueryOneData(DB)
	QueryMultiData(DB)
	//更新数据
	UpdateData(DB)
	// 删除数据
	DeleteData(DB)
}

// CreateTable 创建一个用于测试的表
func CreateTable(DB *sql.DB) {
	sql := `CREATE TABLE IF NOT EXISTS users(
		id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		username VARCHAR(64),
		password VARCHAR(64),
		status INT(4),
		createtime INT(10)
		); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create table failed:", err)
		return
	}
	fmt.Println("create table successd")
}

// InsertData 向数据库里插入一条数据
func InsertData(DB *sql.DB) {
	result, err := DB.Exec("insert INTO users(username,password) values(?,?)", "test", "123456")
	if err != nil {
		fmt.Printf("insert data failed, err: %v", err)
		return
	}
	lastId, err := result.LastInsertId() //本次插入数据的自增id
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("insert succ, id: ", lastId)

	effRows, err := result.RowsAffected() //本次插入操作影响到的行数
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("affected rows", effRows)
}

// QueryOneData 查询一条数据
func QueryOneData(DB *sql.DB) {
	user := new(User1)
	row := DB.QueryRow("select id,username, password, status from users where id=?", 1)
	if err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Status); err != nil {
		fmt.Printf("scan failed, err: %v", err)
		return
	}
	fmt.Println("single row data: ", *user)
}

// QueryMultiData 查询多条数据
func QueryMultiData(DB *sql.DB) {
	user := new(User1)
	rows, err := DB.Query("select id,username,password from users where id = ?", 2)

	defer func() {
		if rows != nil {
			rows.Close() //关闭掉未scan的sql连接
		}
	}()

	if err != nil {
		fmt.Printf("Query failed,err: %v\n", err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Password) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v\n", err)
			return
		}
		fmt.Println("scan successd:", *user)
	}
}

//更新数据
func UpdateData(DB *sql.DB) {
	result, err := DB.Exec("UPDATE users set password=? where id=?", "111111", 1)
	if err != nil {
		fmt.Printf("Insert failed,err:%v\n", err)
		return
	}
	fmt.Println("update data successd:", result)

	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v\n", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}

//删除数据
func DeleteData(DB *sql.DB) {
	result, err := DB.Exec("delete from users where id=?", 1)
	if err != nil {
		fmt.Printf("Insert failed,err:%v\n", err)
		return
	}
	fmt.Println("delete data successd:", result)

	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v\n", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}
