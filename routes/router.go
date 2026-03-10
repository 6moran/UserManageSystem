package routes

import (
	"GoWebUser/controller"
	"GoWebUser/middleware"
	"net/http"
)

func NewRouter(mux *http.ServeMux, uc *controller.UserController) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.Handle("/", middleware.RedirectMiddleware(http.HandlerFunc(uc.RedirectPage)))
	mux.Handle("GET /login", middleware.RedirectMiddleware(http.HandlerFunc(uc.ShowPage)))
	mux.HandleFunc("GET /register", uc.ShowPage)

	mux.Handle("GET /index", middleware.AuthMiddleware(uc.Service)(http.HandlerFunc(uc.ShowPage)))
	mux.Handle("GET /userList", middleware.AuthMiddleware(uc.Service)(http.HandlerFunc(uc.ShowPage)))

	mux.HandleFunc("POST /api/register", uc.HandlerRegister)
	mux.HandleFunc("POST /api/login", uc.HandlerLogin)
	mux.Handle("GET /api/users", middleware.AuthMiddleware(uc.Service)(http.HandlerFunc(uc.HandlerGetUsers)))
	mux.Handle("DELETE /api/users/{id}", middleware.AuthMiddleware(uc.Service)(http.HandlerFunc(uc.HandlerDeleteUser)))
	mux.Handle("GET /api/users/raa", middleware.AuthMiddleware(uc.Service)(http.HandlerFunc(uc.HandlerUserRoleAndAvatar)))
	mux.Handle("PUT /api/users/{id}", middleware.AuthMiddleware(uc.Service)(http.HandlerFunc(uc.HandlerEditUser)))
	mux.HandleFunc("POST /logout", uc.DeleteToken)
}
