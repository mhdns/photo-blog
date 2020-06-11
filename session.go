package main

import (
	"net/http"
	"strconv"
)

var sessionDb = map[int]int{}
var sessionCount int

func createSession(uid int) int {
	sessionCount++
	sessionDb[sessionCount] = uid

	return sessionCount
}

func issueCookie(sessionID int, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "SessionID",
		Value:  strconv.Itoa(sessionID),
		MaxAge: 300,
	})
}

func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("SessionID")
	if err != nil {
		return false
	}

	v, err := strconv.Atoi(c.Value)
	if err != nil {
		return false
	}
	_, ok := userDb[sessionDb[v]]

	return ok
}
