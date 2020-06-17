package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	loginCred := map[string]interface{}{
		"name":     "Anas",
		"password": "asdfghjk",
	}

	body, _ := json.Marshal(loginCred)

	req, err := http.NewRequest("POST", "localhost:5000/login", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Cannot login: %v", err)
	}

	rec := httptest.NewRecorder()

	index(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	fmt.Println(res)
}
