package main

import "io"

var sessionDb map[string]string
var imageDb map[string]string

func testing(w io.Writer) {
	io.WriteString(w, "Working!!")
}
