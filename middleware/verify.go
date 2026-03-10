package middleware

import (
	"GoWebUser/services"
	"GoWebUser/utils"
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

// AuthMiddleware 校验token和用户账号实时状态
func AuthMiddleware(s services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//从Cookie中读取token
			token, err := r.Cookie("token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					fail(w, r, http.StatusUnauthorized, "用户未登录或身份过期")
					return
				}
				//其他系统错误
				fail(w, r, http.StatusInternalServerError, "服务器错误，请稍后再试")
				log.Printf("AuthMiddleware Cookie failed,err:%v\n", err)
				return
			}

			claims, err := utils.ParseToken(token.Value)
			if err != nil {
				if errors.Is(err, utils.ErrorsExpiredToken) {
					http.SetCookie(w, &http.Cookie{
						Name:     "token",
						Value:    "",
						Path:     "/",
						Expires:  time.Unix(0, 0),
						HttpOnly: true,
					})
				}
				fail(w, r, http.StatusUnauthorized, "登录身份已过期")
				return
			}

			//查账号状态
			status, err := s.GetUserStatusByID(claims.UserID)
			if err != nil {
				if errors.Is(err, services.ErrorsIdNotExists) {
					http.SetCookie(w, &http.Cookie{
						Name:     "token",
						Value:    "",
						Path:     "/",
						Expires:  time.Unix(0, 0),
						HttpOnly: true,
					})
					fail(w, r, http.StatusUnauthorized, "您的账户不存在或已被删除")
					return
				}
				fail(w, r, http.StatusInternalServerError, "服务器错误，请稍后再试")
				log.Printf("AuthMiddleware GetUserStatusByID failed,err:%v\n", err)
				return
			}
			if status == 0 {
				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Value:    "",
					Path:     "/",
					Expires:  time.Unix(0, 0),
					HttpOnly: true,
				})
				fail(w, r, http.StatusForbidden, "您的账号已被封禁")
				return
			}

			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RedirectMiddleware 页面跳转中间件
func RedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//从Cookie中读取token
		token, err := r.Cookie("token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)
				return
			}
			//其他系统错误
			fail(w, r, http.StatusInternalServerError, "服务器错误，请稍后再试")
			log.Printf("AuthMiddleware Cookie failed,err:%v\n", err)
			return
		}
		if token.Value != "" {
			http.Redirect(w, r, "/index", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})

}

func fail(w http.ResponseWriter, r *http.Request, status int, msg string) {
	isAPI := strings.HasPrefix(r.URL.Path, "/api/")
	if isAPI {
		utils.SendJSON(w, status, msg, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
