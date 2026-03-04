package controller

import (
	"GoWebUser/models/dto"
	"GoWebUser/services"
	"GoWebUser/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserController struct {
	Service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{Service: service}
}

// render 渲染页面
func render(w http.ResponseWriter, page string) error {
	tmpl, err := template.ParseFiles("template/" + page)
	if err != nil {
		return fmt.Errorf("文件资源不存在,err:%w", err)
	}
	err = tmpl.ExecuteTemplate(w, page, nil)
	if err != nil {
		return err
	}
	return nil
}

// ShowPage 通用展示页面路由函数
func (c *UserController) ShowPage(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	err := render(w, path+".html")
	if err != nil {
		log.Println("模版渲染失败,err:", err)
		http.Error(w, "404，页面不存在", http.StatusNotFound)
	}
}

// HandlerRegister 响应注册
func (c *UserController) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, "参数错误", nil)
		return
	}

	//校验字段
	err = utils.ValidateStruct(req)
	var ve validator.ValidationErrors
	if err != nil {
		//参数校验未通过的错误
		if errors.As(err, &ve) {
			utils.SendJSON(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		//其他系统错误
		utils.SendJSON(w, http.StatusInternalServerError, "参数校验失败", nil)
		return
	}

	err = c.Service.RegisterUser(req)
	if err != nil {
		if errors.Is(err, services.ErrorsEmailExists) {
			utils.SendJSON(w, http.StatusBadRequest, "邮箱已存在", nil)
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，业务处理失败", nil)
		return
	}
	utils.SendJSON(w, http.StatusOK, "注册成功，请登录", nil)
}

// HandlerLogin 响应登录
func (c *UserController) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, "参数错误", nil)
		return
	}

	//校验字段
	err = utils.ValidateStruct(req)
	var ve validator.ValidationErrors
	if err != nil {
		//参数校验未通过的错误
		if errors.As(err, &ve) {
			utils.SendJSON(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		//其他系统错误
		utils.SendJSON(w, http.StatusInternalServerError, "参数校验失败", nil)
		return
	}

	tokenString, err := c.Service.LoginUser(req)
	if err != nil {
		if errors.Is(err, services.ErrorsEmailNotExists) {
			utils.SendJSON(w, http.StatusBadRequest, "邮箱不存在，请先注册", nil)
			return
		}
		if errors.Is(err, services.ErrorsWrongPassword) {
			utils.SendJSON(w, http.StatusBadRequest, "密码错误，请重新输入", nil)
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		return
	}

	//登陆成功需要把token写到cookie里
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(2 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	utils.SendJSON(w, http.StatusOK, "登录成功", nil)
}

func (c *UserController) HandlerGetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageStr := query.Get("page")
	sizeStr := query.Get("size")
	status := query.Get("status")
	keyword := query.Get("keyword")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		return
	}

	//执行业务逻辑
	users, total, err := c.Service.GetUsersByLimit(page, size, status, keyword)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		return
	}

	//执行成功
	utils.SendJSON(w, http.StatusOK, "", map[string]interface{}{
		"list":  users,
		"total": total,
	})
}
