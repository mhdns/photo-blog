package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	name     string
	password string
	uid      int
	role     string
}

var userDb = map[int]user{}
var userCount int

func createUser(u user) (user, error) {
	var newUser user

	if _, ok := userDb[u.uid]; ok {
		return newUser, fmt.Errorf("User already exits")
	}

	if len(u.password) < 8 {
		return newUser, fmt.Errorf("Password must be 8 or more characters")
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.MinCost)
	if err != nil {
		return newUser, fmt.Errorf("Internal server error")
	}

	u.password = string(bs)

	newUser = user{
		u.name,
		u.password,
		u.uid,
		u.role,
	}

	return newUser, nil
}
