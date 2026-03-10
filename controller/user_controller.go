package controller

import (
	"GoWebUser/models/dto"
	"GoWebUser/services"
	"GoWebUser/utils"
	"encoding/json"
	"errors"
	"fmt"
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

// RedirectPage 根页面跳转逻辑
func (c *UserController) RedirectPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

// DeleteToken 响应退出销毁Token
func (c *UserController) DeleteToken(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

// HandlerRegister 响应注册
func (c *UserController) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, "参数错误", nil)
		return
	}

	//校验字段
	err = utils.ValidateStruct(req)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = c.Service.RegisterUser(req)
	if err != nil {
		if errors.Is(err, services.ErrorsEmailExists) {
			utils.SendJSON(w, http.StatusBadRequest, "邮箱已存在", nil)
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerRegistger RegisterUser failed,err:%v\n", err)
		return
	}
	utils.SendJSON(w, http.StatusOK, "注册成功，请登录", nil)
}

// HandlerLogin 响应登录
func (c *UserController) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, "参数错误", nil)
		return
	}

	//校验字段
	err = utils.ValidateStruct(req)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, err.Error(), nil)
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
		if errors.Is(err, services.ErrorsStatusNotMatch) {
			utils.SendJSON(w, http.StatusForbidden, "你的账号已被封禁", nil)
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerLogin LoginUser failed,err:%v\n", err)
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

// HandlerGetUsers 响应请求全部用户信息
func (c *UserController) HandlerGetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageStr := query.Get("page")
	sizeStr := query.Get("size")
	status := query.Get("status")
	keyword := query.Get("keyword")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerGetUsers Atoi1 failed,err:%v\n", err)
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerGetUsers Atoi2 failed,err:%v\n", err)
		return
	}

	//执行业务逻辑
	users, total, err := c.Service.GetUsersByLimit(page, size, status, keyword)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerGetUsers GetUsersByLimit failed,err:%v\n", err)
		return
	}

	//执行成功
	utils.SendJSON(w, http.StatusOK, "", map[string]interface{}{
		"list":  users,
		"total": total,
	})
}

// HandlerDeleteUser 响应删除用户
func (c *UserController) HandlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerDeleteUser Atoi failed,err:%v\n", err)
		return
	}
	err = c.Service.DeleteUserByID(id)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerDeleteUser DeleteUserByID failed,err:%v\n", err)
		return
	}
	utils.SendJSON(w, http.StatusOK, "删除成功", nil)
}

// HandlerUserRoleAndAvatar 响应请求当前用户的角色和头像
func (c *UserController) HandlerUserRoleAndAvatar(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("userID")
	user, err := c.Service.GetUserRoleAndAvatar(id.(int))
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerUserRoleAndAvatar GetUserRoleAndAvatar failed,err:%v\n", err)
		return
	}
	utils.SendJSON(w, http.StatusOK, "", user)
}

// HandlerEditUser 响应修改用户信息
func (c *UserController) HandlerEditUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerEditUser Atoi failed,err:%v\n", err)
		return
	}

	r.ParseMultipartForm(10 << 20)
	username := r.FormValue("username")
	password := r.FormValue("password")
	status, _ := strconv.Atoi(r.FormValue("status"))

	req := dto.EditRequest{
		ID:       id,
		Username: username,
		Password: password,
		Status:   status,
	}
	//校验字段
	err = utils.ValidateStruct(req)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	file, header, _ := r.FormFile("avatar") // 头像文件

	err = c.Service.UpdateUserByID(req, file, header)
	if err != nil {
		if errors.Is(err, services.ErrorsJustImages) {
			utils.SendJSON(w, http.StatusBadRequest, "只允许上传图片", nil)
			return
		}
		if errors.Is(err, services.ErrorsWrongPassword) {
			utils.SendJSON(w, http.StatusBadRequest, "图片大小不能超过2MB", nil)
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, "服务器错误，请稍后再试", nil)
		log.Printf("HandlerEditUser UpdateUserByID failed,err:%v\n", err)
		return
	}
	utils.SendJSON(w, http.StatusOK, "", nil)
}
