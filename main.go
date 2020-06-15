package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
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
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/register", register)
	http.HandleFunc("/home", home)
	http.HandleFunc("/upload", upload)
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("."))))
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
		tpl.ExecuteTemplate(w, "home.gohtml", imageDb[sessionDb[sessionID]])
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
			return
		}

		// Issue Cookie
		issueCookie(createSession(u.uid), w)

		http.Redirect(w, r, "/home", 303)
		return
	}

	tpl.ExecuteTemplate(w, "login.gohtml", nil)
	return
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
	return
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && alreadyLoggedIn(r) {
		// Parse the multipart form
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// Get the file from form
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer file.Close()

		// fmt.Printf("Uploaded File: %+v\n", header.Filename)
		// fmt.Printf("File Size: %+v\n", header.Size)
		// fmt.Printf("MIME Header: %+v\n", header.Header)

		// Get userid
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

		filename := header.Filename

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("public/user_images/%v", filename), fileBytes, 0644)
		if err != nil {
			fmt.Println(err)
		}

		imageDb[userDb[sessionDb[sessionID]].uid] = append(imageDb[userDb[sessionDb[sessionID]].uid], filename)
		tpl.ExecuteTemplate(w, "upload.gohtml", nil)
		return
	}
	if alreadyLoggedIn(r) {
		tpl.ExecuteTemplate(w, "upload.gohtml", nil)
		return
	}

	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		c, err := r.Cookie("SessionID")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		c.MaxAge = -1

		http.SetCookie(w, c)
		http.Redirect(w, r, "/", 303)
	}
}
