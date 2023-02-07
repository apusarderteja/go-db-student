package handler

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/justinas/nosurf"
)

type LoginFormError struct {
	Username string
	Password string
}

type LoginUser struct {
	ID        int          `db:"id" json:"id"`
	FirstName string       `db:"first_name" json:"first_name"`
	LastName  string       `db:"last_name" json:"last_name"`
	Username  string       `db:"username" json:"username"`
	Email     string       `db:"email" json:"email"`
	Password  string       `db:"password" json:"password"`
	Status    bool         `db:"status" json:"status"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}
type LoginUsertemplates struct {
	LoginUser LoginUser
	FormError LoginFormError
	CSRFToken string
}

func (L *LoginUser) Validate() error {
	return validation.ValidateStruct(L,
		validation.Field(&L.Username, validation.Required.Error("This Username cannot be Username")),
		validation.Field(&L.Password, validation.Required.Error("This Filed cannot be Password")),
	)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	pareseLoginTemplate(w, LoginUsertemplates{
		CSRFToken: nosurf.Token(r),
		LoginUser: LoginUser{},
	})
}

func (h Handler) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("%#v", err)
	}

	var usera LoginUser
	if err := h.decoder.Decode(&usera, r.PostForm); err != nil {
		log.Fatal(err)
	}

	if aErr := usera.Validate(); aErr != nil {
		vErrors, ok := aErr.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[key] = value.Error()
			}
			pareseLoginTemplate(w, LoginUsertemplates{

				FormError: LoginFormError{
					Username: "The Username is required.",
					Password: "The password is required.",
				},
				CSRFToken: nosurf.Token(r),
			})
			return
		}
		http.Error(w, aErr.Error(), http.StatusInternalServerError)
		return
	}

	const getuser = `SELECT * FROM users WHERE username= $1 AND password= $2`
	var loginuser LoginUser
	aerr := h.db.Get(&loginuser, getuser, usera.Username, usera.Password)

	if aerr != nil {
		pareseLoginTemplate(w, LoginUsertemplates{
			CSRFToken: nosurf.Token(r),
			FormError: LoginFormError{
				Username: "The Username/password is wrong.",
				Password: "The Username/password is wrong.",
			},
		})
		return
	}
	h.sessionManager.Put(r.Context(), "username", usera.Username)
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func pareseLoginTemplate(w http.ResponseWriter, data any) {
	t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/login.html")
	if err != nil {
		log.Fatalf("%v", err)
	}
	t.ExecuteTemplate(w, "login.html", data)
}
