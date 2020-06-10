package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

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
	if alreadyLoggedIn(r) {
		c, err := r.Cookie("SessionID")
		if err != nil {
			http.Redirect(w, r, "/login", 303)
			return
		}
		sessionID, err := strconv.Atoi(c.Value)
		if err != nil {
			http.Redirect(w, r, "/login", 303)
			return
		}
		tpl.ExecuteTemplate(w, "home.gohtml", userDb[sessionDb[sessionID]].name)
		return
	}
	http.Redirect(w, r, "/login", 303)
}

func login(w http.ResponseWriter, r *http.Request) {

	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/home", 303)
		return
	} else if r.Method == http.MethodPost {
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

		// Issue Cookie
		issueCookie(createSession(u.uid), w)

		http.Redirect(w, r, "/home", 303)
		return
	}

	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/home", 303)
		return
	} else if r.Method == http.MethodPost {

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

		// Issue Cookie
		issueCookie(createSession(u.uid), w)

		http.Redirect(w, r, "/", 303)
		return
	}
	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "upload.gohtml", nil)
}
