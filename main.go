package main

import (
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/home", home)
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":5000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "html/text")
	tpl.ExecuteTemplate(w, "index.gohtml", nil)

}

func home(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "home.gohtml", nil)
}

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		name := r.PostFormValue("name")
		password := r.PostFormValue("password")

		checkDb := func(name string, db map[int]user) (user, error) {
			for _, v := range db {
				if v.name == name {
					return v, nil
				}
			}
			var u user
			return u, fmt.Errorf("User not found")
		}

		u, err := checkDb(name, userDb)
		if err != nil {
			http.Error(w, "Invalid Credentials", 400)
			return
		}

		// Check password
		err = bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid Credentials", 400)
		}

		http.Redirect(w, r, "/home", 303)
		return
	}

	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		uid := userCount + 1
		role := "user"
		userCount++
		name := r.PostFormValue("name")
		password := r.PostFormValue("password")

		u := user{
			name,
			password,
			uid,
			role,
		}

		user, err := createUser(u)
		if err != nil {
			http.Error(w, err.Error(), 400)
			uid--
			return
		}

		userDb[uid] = user

		fmt.Println(user)
		http.Redirect(w, r, "/", 303)

		fmt.Println(userDb)
		return
	}
	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "upload.gohtml", nil)
}
