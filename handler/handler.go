package handler

import (
	"DB/web/storage/postgres"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/form"
)

type Handler struct {
	sessionManager *scs.SessionManager
	decoder        *form.Decoder
	storage             postgres.PostgresStorage
}

func NewHandler(sm *scs.SessionManager, formDecoder *form.Decoder, db *postgres.PostgresStorage) *chi.Mux {
	h := &Handler{
		sessionManager: sm,
		decoder:        formDecoder,
		db:             *storage,
	}

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(Method)

	r.Group(func(r chi.Router) {
		r.Use(sm.LoadAndSave)
		r.Get("/", h.Home)
		r.Get("/login", h.Login)
		r.Post("/login", h.LoginPostHandler)
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))

	r.Group(func(r chi.Router) {
		r.Use(sm.LoadAndSave)
		r.Use(h.Authentication)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.ListUser)

			r.Get("/create", h.CreateUser)

			r.Post("/store", h.StoreUser)

			r.Get("/{id:[0-9]+}/edit", h.EditUser)

			r.Put("/{id:[0-9]+}/update", h.UpdateUser)

			r.Get("/{id:[0-9]+}/delete", h.DeleteUser)
		})

		r.Get("/logout", h.LogoutHandler)
	})

	return r
}

func Method(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch strings.ToLower(r.PostFormValue("_method")) {
			case "put":
				r.Method = http.MethodPut
			case "patch":
				r.Method = http.MethodPatch
			case "delete":
				r.Method = http.MethodDelete
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h Handler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := h.sessionManager.GetString(r.Context(), "username")
		log.Println("username: ", username)
		if username == "" {
			// http.Error(w, "unauthorized", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
