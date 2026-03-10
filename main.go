package main

import (
	"GoWebUser/controller"
	"GoWebUser/infrastructure"
	"GoWebUser/repositories/repo_mysql"
	"GoWebUser/routes"
	"GoWebUser/services"
	"fmt"
	"log"
	"net/http"
)

func main() {
	//初始化数据库
	db, err := infrastructure.InitDB()
	if err != nil {
		log.Printf("连接数据库失败,err:%v\n", err)
	}
	fmt.Println("连接数据库成功")

	//依赖注入
	userRepo := repo_mysql.NewMySQLUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	//自定义路由
	mux := http.NewServeMux()
	//注册路由
	routes.NewRouter(mux, userController)

	addr := "0.0.0.0:8080"
	fmt.Printf("服务器正在启动,监听地址为:%v\n", addr)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Printf("服务器启动失败,err:%v\n", err)
	}

}
