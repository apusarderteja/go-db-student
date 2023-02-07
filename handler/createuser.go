package handler

import (
	// "fmt"
	"fmt"
	"log"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/justinas/nosurf"
)

func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	pareseCreateUserTemplate(w, User{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) StoreUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("%#v", err)
	}

	user := User{}
	if err := h.decoder.Decode(&user, r.PostForm); err != nil {
		log.Fatal(err)
		fmt.Println(user)
	}
	if err := user.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			for key, val := range vErr {
				user.FormError[strings.Title(key)] = val
			}
		}
		pareseCreateUserTemplate(w, user)
		return
	}
	const insertQuery = `
		INSERT INTO users(
			first_name,
			last_name,
			username,
			email,
			password
		) VALUES (
			:first_name,
			:last_name,
			:username,
			:email,
			:password
		) RETURNING id;
	`

	stmt, err := h.db.PrepareNamed(insertQuery)
	if err != nil {
		log.Fatalln(err)
	}
	var uID int
	err = stmt.Get(&uID, user)
	if err != nil {
		log.Fatal(err)
	}

	if uID == 0 {
		log.Fatalln("unable to insert user")
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
