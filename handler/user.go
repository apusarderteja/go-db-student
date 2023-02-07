package handler

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID        int              `json:"id" form:"-" db:"id"`
	FirstName string           `json:"first_name" db:"first_name"`
	LastName  string           `json:"last_name" db:"last_name"`
	Email     string           `json:"email" db:"email"`
	Username  string           `json:"username" db:"username"`
	Password  string           `json:"password" db:"password"`
	Status    bool             `json:"status" db:"status"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime     `json:"deleted_at" db:"deleted_at"`
	FormError map[string]error `json:"-" form:"-"`
	CSRFToken string           `json:"-" form:"csrf_token"`
}

type UserList struct {
	Users []User `json:"users"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName,
			validation.Required.Error("The first name field is required."),
			validation.Length(3, 32).Error("The first name field must be between 3 to 32 characters."),
		),
		validation.Field(&u.LastName,
			validation.Required.Error("The last name field is required."),
			validation.Length(3, 32).Error("The last name field must be between 3 to 32 characters."),

		),
		validation.Field(&u.Username,
			validation.Required.When(u.ID==0).Error("The username field is required."),
		),
		validation.Field(&u.Email,
			validation.Required.When(u.ID==0).Error("The email field is required."),
			is.Email.Error("The email field must be a valid email."),
		),
		validation.Field(&u.Password,
			validation.Required.When(u.ID==0).Error("The password field is required."),
		),
	)
}

func (h Handler) ListUser(w http.ResponseWriter, r *http.Request) {
	const listQuery = `SELECT * from users WHERE deleted_at IS NULL ORDER BY id ASC`
	var listUser []User
	if err := h.db.Select(&listUser, listQuery); err != nil {
		log.Fatalln(err)
	}

	t, err := template.ParseFiles("templates/admin/user/list-user.html")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("%+v", listUser)
	t.Execute(w, listUser)
}

func pareseCreateUserTemplate(w http.ResponseWriter, data any) {
	t := template.New("create user")
	t = template.Must(t.ParseFiles("templates/admin/user/create-user.html", "templates/admin/user/_form.html"))

	if err := t.ExecuteTemplate(w, "create-user.html", data); err != nil {
		log.Fatal(err)
	}
}
