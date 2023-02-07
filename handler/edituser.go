package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	// "strings"

	// "strings"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println(id)
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	const getuser = `SELECT * FROM users WHERE id=$1 AND deleted_at IS NULL`
	var user User
	user.ID = uID
	if err := h.db.Get(&user, getuser, uID) ; err!= nil {
		log.Fatalln(err)
	}

	user.CSRFToken = nosurf.Token(r)
	pareseEditUserTemplate(w, user)
}

func pareseEditUserTemplate(w http.ResponseWriter, data any) {
	var err error
	t := template.New("edit user")
	t = template.Must(t.ParseFiles("templates/header.html", "templates/footer.html", "templates/admin/user/edit-user.html","templates/admin/user/_form.html"))
	if err != nil {
		log.Fatalf("%v", err)
	}

	t.ExecuteTemplate(w, "edit-user.html", data)
}

func (h Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	user := User{}
	if err := h.decoder.Decode(&user, r.PostForm); err != nil {
		log.Fatal(err)
	}

	user.ID = uID

	if err := user.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
				user.FormError = vErr
		}
		fmt.Println(user)
		pareseEditUserTemplate(w, user)
		return
	}

	const updateQuery = `UPDATE users
		SET first_name = :first_name, 
		last_name = :last_name,
		status = :status
		WHERE id = :id
		RETURNING id;
	`

	stmt, err := h.db.PrepareNamed(updateQuery)
	if err != nil {
		log.Fatal(err)
	}

	if err := stmt.Get(&uID, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)

}