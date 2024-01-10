package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var (
	port int = 8081
	addr string
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	var tmplFile = "index.html"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var tmplFile = "login.html"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func init() {
	addr = fmt.Sprintf("127.0.0.1:%v", port)
}

func main() {
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/login", LoginHandler)

	fmt.Printf("Starting server on %v.\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
