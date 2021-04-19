package server

import (
	"net/http"
	"github.com/Notarrogantjustbetter/AuthenticationSystem/v2/database"
	"github.com/Notarrogantjustbetter/AuthenticationSystem/v2/middleware"
	"github.com/Notarrogantjustbetter/AuthenticationSystem/v2/sessions"
	"github.com/Notarrogantjustbetter/AuthenticationSystem/v2/utils"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func (s Server) InitServer() *mux.Router {
	goServer := &Server{
		router: mux.NewRouter(),
	}
	goServer.routesHandler()
	return goServer.router
}

func (s Server) routesHandler() {
	s.router.HandleFunc("/", middleware.MiddlewareAuthentication(s.homeHandler().ServeHTTP)).Methods("GET")
	s.router.HandleFunc("/register", s.registerHandler().ServeHTTP).Methods("GET", "POST")
	s.router.HandleFunc("/login", s.loginHandler().ServeHTTP).Methods("GET", "POST")
	s.router.HandleFunc("/logout", s.logoutHandler().ServeHTTP).Methods("GET", "POST")
}

func (s Server) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.ExecuteTemplate(w, "home.html", nil)
	}
}

func (s Server) registerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			utils.ExecuteTemplate(w, "register.html", nil)
		case "POST":
			r.ParseForm()
			nickname := r.PostForm.Get("Nickname")
			password := r.PostForm.Get("Password")
			err := database.RegisterUser(nickname, password)
			if err == database.ErrUsernameTaken {
				utils.ExecuteTemplate(w, "register.html", "username is taken")
				return
			} else if err != nil {
				utils.InternalServerError(w)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

func (s Server) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			utils.ExecuteTemplate(w, "login.html", nil)
		case "POST":
			r.ParseForm()
			nickname := r.PostForm.Get("Nickname")
			password := r.PostForm.Get("Password")
			user, err := database.LoginUser(nickname, password)
			if err != nil {
				switch err {
				case database.ErrUserNotFound:
					utils.ExecuteTemplate(w, "login.html", "user not found")
				case database.ErrInvalidLogin:
					utils.ExecuteTemplate(w, "login.html", "invalid login")
				default:
					utils.InternalServerError(w)
				}
				return
			}
			userId, err := user.GetId()
			if err != nil {
				utils.InternalServerError(w)
				return
			}
			sessions.SetSession(w, r, userId)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func (s Server) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions.DeleteSession(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}