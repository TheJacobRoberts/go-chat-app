package main

import (
	"fmt"
	"net/http"
)

var (
	port int = 8080
	addr string
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is my content.")
	fmt.Fprintln(w, r.Header)
	fmt.Fprintln(w, r.Body)
}

func init() {
	addr = fmt.Sprintf("127.0.0.1:%v", port)
}

func main() {
	http.HandleFunc("/", RootHandler)

	fmt.Printf("Starting server on %v.\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
