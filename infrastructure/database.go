package infrastructure

import (
	"database/sql"
	"fmt"
)

func InitDB() (*sql.DB, error) {
	//数据库信息
	dsn := "root:1458963@tcp(127.0.0.1:3306)/user_management_system?parseTime=true"
	//连接数据库
	db, err := sql.Open("mysql", dsn)
	//不会检验用户名和密码是否正确
	//dsn格式不对的时候报这个错
	if err != nil {
		return nil, fmt.Errorf("dsn格式错误,err:%w", err)
	}

	//尝试ping，检验用户名和密码是否正确
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误,err:%w", err)
	}
	return db, nil
}
