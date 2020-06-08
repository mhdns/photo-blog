package main

import (
	"fmt"
	"html/template"
	"net/http"
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
