package db

import (
	mydb "DistributedMemory/db/mysql"
	"fmt"
)

// UserSignup 通过用户名及密码完成user表的注册操作
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		fmt.Println("Failed to insert, err:", err)
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:", err)
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}

	return false
}

// UserSignin 判断密码是否一致
func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select user_pwd from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	var pwd string
	err = stmt.QueryRow(username).Scan(&pwd)
	if err != nil {
		fmt.Println("username not found")
		return false
	}
	if pwd == encpwd {
		return true
	}

	return false
}
